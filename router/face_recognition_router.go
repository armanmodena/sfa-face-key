package router

import (
	"arkan-face-key/handler"
	"arkan-face-key/service"

	"github.com/gin-gonic/gin"
	"github.com/pkg/sftp"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetupFaceRecognitionRouter(r *gin.Engine, mongo *mongo.Client, sftp *sftp.Client) {
	sftpService := service.NewSftpService(sftp)
	faceService := service.NewFaceRecognitionService(mongo, sftpService)
	faceHandler := handler.NewFaceRecognitionHandler(faceService)

	api := r.Group("/api")
	{
		api.POST("/face/save", faceHandler.SaveUserFaceKey)
		api.POST("/face/validate/embedding", faceHandler.ValidateWithEmbedding)
		api.POST("/face/validate/image", faceHandler.ValidateWithImage)
	}
}
