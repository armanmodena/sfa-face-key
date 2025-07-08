package service

import (
	"arkan-face-key/config"
	"arkan-face-key/helper"
	"arkan-face-key/model"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/Kagami/go-face"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type FaceRecognitionService interface {
	SaveUserFaceKey(r *gin.Context, image *multipart.FileHeader, username string) (*helper.Response, *helper.Response)
	ValidateWithEmbedding(r *gin.Context, image *multipart.FileHeader, username string, threshold float32) (*helper.Response, *helper.Response)
	ValidateWithImage(r *gin.Context, image *multipart.FileHeader, username string, threshold float32) (*helper.Response, *helper.Response)
}

type faceRecognitionService struct {
	mongo       *mongo.Client
	sftpService SftpService
}

func NewFaceRecognitionService(mongo *mongo.Client, sftpService SftpService) FaceRecognitionService {
	return &faceRecognitionService{mongo: mongo, sftpService: sftpService}
}

const dataDir = "faces"

func (s *faceRecognitionService) SaveUserFaceKey(r *gin.Context, image *multipart.FileHeader, username string) (*helper.Response, *helper.Response) {
	// Validate username
	if username == "" {
		return nil, &helper.Response{
			Status:  400,
			Message: "Invalid Username",
		}
	}

	// Get user from database
	var user model.User
	collection := s.mongo.Database(config.MONGO_DB).Collection("user")
	err := collection.FindOne(r, map[string]any{"username": username}).Decode(&user)
	if err != nil {
		return nil, &helper.Response{
			Status:  404,
			Message: fmt.Sprintf("User with username %s not found: %v", username, err),
		}
	}

	// Initialize the face recognizer
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Can't init face recognizer: %v", err),
		}
	}
	defer rec.Close()

	// Open the uploaded image file
	file, err := image.Open()
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error opening uploaded image: %v", err),
		}
	}
	defer file.Close()

	// Read the uploaded file into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error reading file byte uploaded image: %v", err),
		}
	}

	// Recognize faces in the uploaded image
	refFace, err := rec.Recognize(fileBytes)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error recognizing face in uploaded image: %v", err),
		}
	}

	// Check if any faces were found
	if len(refFace) == 0 {
		return nil, &helper.Response{
			Status:  400,
			Message: "No faces found in the image",
		}
	}

	// Check if multiple faces were found
	if len(refFace) > 1 {
		return nil, &helper.Response{
			Status:  400,
			Message: "Multiple faces found in the image",
		}
	}

	// Extract descriptors (embeddings)
	embedding := refFace[0].Descriptor
	faceKeyFileName := fmt.Sprintf("%s_%d_face_key.jpeg", user.Username, time.Now().Unix())
	// Upload the file to SFTP
	// Upload the file to SFTP
	_, errRes := s.sftpService.UploadFile(image, faceKeyFileName)
	if errRes != nil {
		return nil, &helper.Response{
			Status:  errRes.Status,
			Message: fmt.Sprintf("Error uploading face key file: %v", errRes.Message),
		}
	}

	// Convert embedding to string for storage
	embeddingStr, err := helper.DescriptorToString(embedding)
	if err != nil {
		return nil, &helper.Response{
			Status:  500,
			Message: fmt.Sprintf("Error converting embedding to string: %v", err),
		}
	}

	// Delete old face key file if it exists
	if user.GoFaceImageUrl != "" {
		_, errRes := s.sftpService.DeleteFile(user.GoFaceImageUrl)
		if errRes != nil {
			return nil, &helper.Response{
				Status:  errRes.Status,
				Message: fmt.Sprintf("Error deleting old face key file: %v", errRes.Message),
			}
		}
	}

	//save to database
	update := map[string]any{
		"go_face_image_url": faceKeyFileName,
		"go_face_embedding": embeddingStr,
	}
	_, err = collection.UpdateOne(
		r,
		map[string]any{"username": user.Username},
		map[string]any{"$set": update},
	)
	if err != nil {
		return nil, &helper.Response{
			Status:  500,
			Message: "Error saving user embedding",
			Data: map[string]any{
				"error": fmt.Sprintf("Error saving user embedding: %v", err),
			},
		}
	}

	return &helper.Response{
		Status:  200,
		Message: "Face key saved successfully",
		Data: map[string]any{
			"user_id":            user.Id,
			"face_key_file":      faceKeyFileName,
			"face_key_embedding": embeddingStr,
		},
	}, nil
}

func (s *faceRecognitionService) ValidateWithEmbedding(r *gin.Context, image *multipart.FileHeader, username string, threshold float32) (*helper.Response, *helper.Response) {
	// Validate username
	if username == "" {
		return nil, &helper.Response{
			Status:  400,
			Message: "Invalid Username",
		}
	}

	// Get user from database
	var user model.User
	collection := s.mongo.Database(config.MONGO_DB).Collection("user")
	err := collection.FindOne(r, map[string]any{"username": username}).Decode(&user)
	if err != nil {
		return nil, &helper.Response{
			Status:  404,
			Message: fmt.Sprintf("User with username %s not found: %v", username, err),
		}
	}

	// Check if user has a face key embedding
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Can't init face recognizer: %v", err),
		}
	}
	defer rec.Close()

	// Check if user has a face key file
	file, err := image.Open()
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error opening uploaded image: %v", err),
		}
	}
	defer file.Close()

	// Read the uploaded file into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error reading uploaded image: %v", err),
		}
	}

	// Recognize faces in the uploaded image
	refFace, err := rec.Recognize(fileBytes)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error recognizing face in uploaded image: %v", err),
		}
	}

	// Check if any faces were found
	if len(refFace) == 0 {
		return nil, &helper.Response{
			Status:  400,
			Message: "No faces found",
		}
	}

	// Check if multiple faces were found
	if len(refFace) > 1 {
		return nil, &helper.Response{
			Status:  400,
			Message: "Multiple faces found",
		}
	}

	// Extract descriptors (embeddings)
	desc1 := refFace[0].Descriptor
	// Convert user.FaceKeyEmbedding (JSON string) to face.Descriptor
	desc2, err := helper.StringToDescriptor(user.GoFaceEmbedding)
	if err != nil {
		return nil, &helper.Response{
			Status:  500,
			Message: fmt.Sprintf("Error converting embedding string to descriptor: %v", err),
		}
	}

	// Compute Euclidean distance
	distance := euclideanDistance(desc1, desc2)

	// Check if the distance is below the threshold
	if distance > threshold {
		return nil, &helper.Response{
			Status:  400,
			Message: "Face not matched",
		}
	}

	return &helper.Response{
		Status:  200,
		Message: "Face matched",
		Data: map[string]any{
			"user_id":            user.Id,
			"username":           user.Username,
			"full_name":          user.FullName,
			"face_key_file":      user.GoFaceImageUrl,
			"face_key_embedding": user.GoFaceEmbedding,
		},
	}, nil
}

func (s *faceRecognitionService) ValidateWithImage(r *gin.Context, image *multipart.FileHeader, username string, threshold float32) (*helper.Response, *helper.Response) {
	// Validate username
	if username == "" {
		return nil, &helper.Response{
			Status:  400,
			Message: "Invalid Username",
		}
	}

	// Get user from database
	var user model.User
	collection := s.mongo.Database(config.MONGO_DB).Collection("user")
	err := collection.FindOne(r, map[string]any{"username": username}).Decode(&user)
	if err != nil {
		return nil, &helper.Response{
			Status:  404,
			Message: fmt.Sprintf("User with username %s not found: %v", username, err),
		}
	}

	// Initialize the face recognizer
	rec, err := face.NewRecognizer(dataDir)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Can't init face recognizer: %v", err),
		}
	}
	defer rec.Close()

	// Check if user has a face key file
	if user.GoFaceImageUrl == "" {
		return nil, &helper.Response{
			Status:  400,
			Message: "User does not have a face key image",
		}
	}

	// Load the base image for the user
	baseImage := filepath.Join(dataDir, "images/"+user.GoFaceImageUrl)
	baseFaces, err := rec.RecognizeFile(baseImage)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error recognizing face in base image: %v", err),
		}
	}

	// Check if any faces were found in the base image
	if len(baseFaces) == 0 {
		return nil, &helper.Response{
			Status:  400,
			Message: "No faces found",
		}
	}

	// Check if multiple faces were found in the base image
	if len(baseFaces) > 1 {
		return nil, &helper.Response{
			Status:  400,
			Message: "Multiple faces found in the base image",
		}
	}

	// Use the first face found in the base image
	baseFace := baseFaces[0]

	// Prepare the recognizer with the base face descriptor
	var samples []face.Descriptor
	var sampleIndexes []int32

	// Add the base face descriptor to the recognizer
	samples = append(samples, baseFace.Descriptor)
	sampleIndexes = append(sampleIndexes, int32(0))

	// Set the samples and their indexes in the recognizer
	rec.SetSamples(samples, sampleIndexes)

	// Open the uploaded image file
	file, err := image.Open()
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error opening uploaded image: %v", err),
		}
	}
	defer file.Close()

	// Read the uploaded file into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error reading uploaded image: %v", err),
		}
	}

	// Recognize faces in the uploaded image
	faces, err := rec.Recognize(fileBytes)
	if err != nil {
		return nil, &helper.Response{
			Status:  400,
			Message: fmt.Sprintf("Error recognizing face in uploaded image: %v", err),
		}
	}

	// Check if any faces were found in the uploaded image
	if len(faces) == 0 {
		return nil, &helper.Response{
			Status:  400,
			Message: "No faces found",
		}
	}

	// Check if multiple faces were found in the uploaded image
	if len(faces) > 1 {
		return nil, &helper.Response{
			Status:  400,
			Message: "Multiple faces found in the uploaded image",
		}
	}

	// Classify the face in the uploaded image
	faceIndex := rec.ClassifyThreshold(faces[0].Descriptor, threshold)
	if faceIndex < 0 {
		return nil, &helper.Response{
			Status:  400,
			Message: "Face not matched",
		}
	}

	// Check if the recognized face matches the user's face
	return &helper.Response{
		Status:  200,
		Message: "Face matched",
		Data: map[string]any{
			"user_id":            user.Id,
			"username":           user.Username,
			"full_name":          user.FullName,
			"face_key_file":      user.GoFaceImageUrl,
			"face_key_embedding": user.GoFaceEmbedding,
		},
	}, nil
}

// euclideanDistance calculates the Euclidean distance between two face descriptors.
func euclideanDistance(a, b face.Descriptor) float32 {
	var sum float32
	for i := range a {
		diff := a[i] - b[i]
		sum += diff * diff
	}
	return float32((sum) * 0.5) // Remove *0.5 if you want the actual Euclidean distance (sqrt(sum))
}
