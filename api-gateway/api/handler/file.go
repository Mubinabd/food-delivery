package handler

import (
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"
)

// @Summary Upload a file to MinIO
// @Description Upload a file to MinIO
// @Tags MinIO
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param filename formData string true "Filename"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /minio/upload [post]
func (h *HandlerStruct) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	filename := c.PostForm("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	// Extract file extension
	ext := filepath.Ext(header.Filename)
	finalFilename := filename + ext

	bucketName := "testbucket"
	objectName := "images/" + time.Now().Format("20060102150405") + "-" + finalFilename

	// Ensure the bucket exists
	exists, err := h.Clients.MinIOClient.BucketExists(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check if bucket exists", "details": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bucket does not exist"})
		return
	}

	// Upload the file to MinIO
	_, err = h.Clients.MinIOClient.PutObject(bucketName, objectName, file, header.Size, minio.PutObjectOptions{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image", "details": err.Error()})
		return
	}

	// Generate a presigned URL
	presignedURL, err := h.Clients.MinIOClient.PresignedGetObject(bucketName, objectName, 24*time.Hour, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": presignedURL.String()})
}
