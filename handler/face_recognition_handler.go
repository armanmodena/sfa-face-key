package handler

import (
	"arkan-face-key/helper"
	"arkan-face-key/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FaceRecognitionHandler struct {
	service service.FaceRecognitionService
}

func NewFaceRecognitionHandler(service service.FaceRecognitionService) *FaceRecognitionHandler {
	return &FaceRecognitionHandler{service}
}

func (h *FaceRecognitionHandler) SaveUserFaceKey(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.Response{
			Status:  400,
			Message: "Image file is required",
		})
		return
	}

	// userId from form-data as int
	username := c.PostForm("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, helper.Response{
			Status:  400,
			Message: "Username is required",
		})
		return
	}

	res, errRes := h.service.SaveUserFaceKey(c, image, username)
	if errRes != nil {
		c.JSON(errRes.Status, helper.Response{
			Status:  errRes.Status,
			Message: errRes.Message,
		})
		return
	}

	c.JSON(http.StatusOK, helper.Response{
		Status:  res.Status,
		Message: res.Message,
		Data:    res.Data,
	})
}

func (h *FaceRecognitionHandler) ValidateWithEmbedding(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.Response{
			Status:  400,
			Message: "Image file is required",
		})
		return
	}

	// userId from form-data as int
	username := c.PostForm("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, helper.Response{
			Status:  400,
			Message: "Username is required",
		})
		return
	}

	thresholdStr := c.PostForm("threshold")
	var threshold float32
	if thresholdStr != "" {
		var threshold64 float64
		threshold64, err = strconv.ParseFloat(thresholdStr, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.Response{
				Status:  400,
				Message: "Threshold must be a valid float",
			})
			return
		}
		threshold = float32(threshold64)
	} else {
		// Default threshold if not provided
		threshold = 0.6
	}

	res, errRes := h.service.ValidateWithEmbedding(c, image, username, threshold)
	if errRes != nil {
		c.JSON(errRes.Status, helper.Response{
			Status:  errRes.Status,
			Message: errRes.Message,
		})
		return
	}

	c.JSON(http.StatusOK, helper.Response{
		Status:  res.Status,
		Message: res.Message,
		Data:    res.Data,
	})
}

func (h *FaceRecognitionHandler) ValidateWithImage(c *gin.Context) {
	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, helper.Response{
			Status:  400,
			Message: "Image file is required",
		})
		return
	}

	// userId from form-data as int
	username := c.PostForm("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, helper.Response{
			Status:  400,
			Message: "User ID is required",
		})
		return
	}

	thresholdStr := c.PostForm("threshold")
	var threshold float32
	if thresholdStr != "" {
		var threshold64 float64
		threshold64, err = strconv.ParseFloat(thresholdStr, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, helper.Response{
				Status:  400,
				Message: "Threshold must be a valid float",
			})
			return
		}
		threshold = float32(threshold64)
	} else {
		// Default threshold if not provided
		threshold = 0.6
	}

	res, errRes := h.service.ValidateWithImage(c, image, username, threshold)
	if errRes != nil {
		c.JSON(errRes.Status, helper.Response{
			Status:  errRes.Status,
			Message: errRes.Message,
		})
		return
	}

	c.JSON(http.StatusOK, helper.Response{
		Status:  res.Status,
		Message: res.Message,
	})
}
