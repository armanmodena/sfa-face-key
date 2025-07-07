package helper

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/Kagami/go-face"
)

type Response struct {
	Meta    any    `json:"meta,omitempty"`
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// add function descriptorToString
func DescriptorToString(descriptor face.Descriptor) (string, error) {
	jsonData, err := json.Marshal(descriptor)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// add function stringToDescriptor
func StringToDescriptor(data string) (face.Descriptor, error) {
	var descriptor face.Descriptor
	err := json.Unmarshal([]byte(data), &descriptor)
	if err != nil {
		return face.Descriptor{}, err
	}
	return descriptor, nil
}

// normalizeExtension maps file extensions to a consistent format
func normalizeExtension(ext string) string {
	ext = strings.TrimPrefix(ext, ".")
	ext = strings.ToLower(ext)

	switch ext {
	case "jpg":
		return "jpeg"
	case "jpeg":
		return "jpeg"
	case "png":
		return "png"
	case "webp":
		return "webp"
	case "bmp":
		return "bmp"
	case "gif":
		return "gif"
	default:
		return ext // fallback to the raw normalized extension
	}
}

// GetFileExtensionFromHeader returns a consistent extension (e.g., "jpeg") from multipart.FileHeader
func GetFileExtensionFromHeader(header *multipart.FileHeader) string {
	ext := filepath.Ext(header.Filename)
	fmt.Println("File extension from header:", ext)
	return normalizeExtension(ext)
}

// GetFileExtensionFromPath returns a consistent extension (e.g., "jpeg") from a file path
func GetFileExtensionFromPath(path string) string {
	ext := filepath.Ext(path)
	fmt.Println("File extension from path:", ext)
	return normalizeExtension(ext)
}
