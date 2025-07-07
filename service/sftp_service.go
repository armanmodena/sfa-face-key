package service

import (
	"arkan-face-key/config"
	"arkan-face-key/helper"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/pkg/sftp"
)

type SftpService interface {
	UploadFile(file *multipart.FileHeader, fileName string) (*helper.Response, *helper.Response)
	DeleteFile(fileName string) (*helper.Response, *helper.Response)
	DownloadFile(fileName string) (*helper.Response, *helper.Response)
	GetListOfFile() (*helper.Response, *helper.Response)
}

type sftpService struct {
	sftp *sftp.Client
}

func NewSftpService(sftp *sftp.Client) SftpService {
	return &sftpService{sftp: sftp}
}

func (s *sftpService) UploadFile(file *multipart.FileHeader, fileName string) (*helper.Response, *helper.Response) {
	if s.sftp == nil {
		return nil, &helper.Response{
			Status:  http.StatusInternalServerError,
			Message: "SFTP client is not initialized",
		}
	}

	srcFile, err := file.Open()
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	defer srcFile.Close()

	dstPath := config.SFTP_ROOT + "face_key/" + fileName
	dstFile, err := s.sftp.Create(dstPath)
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}
	defer dstFile.Close()

	_, err = dstFile.ReadFrom(srcFile)
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return &helper.Response{
		Status:  http.StatusOK,
		Message: "File uploaded successfully",
		Data:    dstPath,
	}, nil
}

func (s *sftpService) DeleteFile(fileName string) (*helper.Response, *helper.Response) {
	if s.sftp == nil {
		return nil, &helper.Response{
			Status:  http.StatusInternalServerError,
			Message: "SFTP client is not initialized",
		}
	}

	dstPath := config.SFTP_ROOT + "face_key/" + fileName
	err := s.sftp.Remove(dstPath)
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return &helper.Response{
		Status:  http.StatusOK,
		Message: "File deleted successfully",
	}, nil
}

func (s *sftpService) DownloadFile(fileName string) (*helper.Response, *helper.Response) {
	if s.sftp == nil {
		return nil, &helper.Response{
			Status:  http.StatusInternalServerError,
			Message: "SFTP client is not initialized",
		}
	}

	dstPath := config.SFTP_ROOT + "face_key/" + fileName
	srcFile, err := s.sftp.Open(dstPath)
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	tmpFilePath := "tmp_file/" + fileName
	localFile, err := os.Create(tmpFilePath)
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	defer localFile.Close()

	_, err = srcFile.WriteTo(localFile)
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	defer srcFile.Close()

	return &helper.Response{
		Status:  http.StatusOK,
		Message: "File downloaded and saved to tmp_file successfully",
		Data:    tmpFilePath,
	}, nil
}

func (s *sftpService) GetListOfFile() (*helper.Response, *helper.Response) {
	if s.sftp == nil {
		return nil, &helper.Response{
			Status:  http.StatusInternalServerError,
			Message: "SFTP client is not initialized",
		}
	}

	files, err := s.sftp.ReadDir(config.SFTP_ROOT + "face_key/")
	if err != nil {
		return nil, &helper.Response{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		}
	}

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	return &helper.Response{
		Status:  http.StatusOK,
		Message: "List of files retrieved successfully",
		Data:    fileNames,
	}, nil
}
