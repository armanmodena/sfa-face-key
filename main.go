package main

import (
	"arkan-face-key/config"
	"arkan-face-key/middleware"
	"arkan-face-key/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.InitTimeZone()
	if err != nil {
		log.Fatalf("Failed to load timezone: %v", err)
	}

	mdb, err := config.OpenMongoConnection()
	if err != nil {
		log.Fatal("Error connecting to MongoDB")
	}

	sftp, err := config.OpenSFTPConnection()
	if err != nil {
		log.Fatal("Error connecting to SFTP server")
	}

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware())

	router.SetupFaceRecognitionRouter(r, mdb, sftp)

	port := config.PORT
	if port == "" {
		port = "9050"
	}
	log.Println("Starting server on port " + port)
	r.Run(":" + port)
}
