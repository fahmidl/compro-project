package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	S3Client          *s3.Client
	Bucket            string
	CloudFrontDomain  string
}

func NewUploadHandler(s3Client *s3.Client) *UploadHandler {
	return &UploadHandler{
		S3Client:         s3Client,
		Bucket:           os.Getenv("S3_UPLOADS_BUCKET"),
		CloudFrontDomain: os.Getenv("CLOUDFRONT_DOMAIN"),
	}
}

func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no image file provided"})
		return
	}
	defer file.Close()

	// Validate file type
	ext := filepath.Ext(header.Filename)
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type, allowed: jpg, jpeg, png, gif, webp"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	s3Key := fmt.Sprintf("uploads/%s", filename)

	// Determine content type
	contentType := "application/octet-stream"
	switch ext {
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// Upload to S3
	_, err = h.S3Client.PutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.Bucket),
		Key:         aws.String(s3Key),
		Body:        io.Reader(file),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file to S3"})
		return
	}

	// Return the CloudFront URL if configured, otherwise return the S3 path
	var imageURL string
	if h.CloudFrontDomain != "" {
		imageURL = fmt.Sprintf("/uploads/%s", filename)
	} else {
		imageURL = fmt.Sprintf("/uploads/%s", filename)
	}
	c.JSON(http.StatusOK, gin.H{"url": imageURL})
}
