package apihandlers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"r2-api-go/cache"
	"r2-api-go/config"
	"r2-api-go/r2"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func generateDownloadURL(objectKey string) (string, error) {
	downloadURL, err := r2.GenerateDownloadURL(objectKey)
	if err != nil {
		return "", errors.New("failed to generate download URL")
	}
	return *downloadURL, nil
}

func UploadImage(c *gin.Context) {
	tempFilename := c.Query("fileName")
	if tempFilename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Filename is required"})
		return
	}

	fmt.Println("Uploaded Filename:", tempFilename)
	tempFilePath := filepath.Join("uploads", tempFilename)

	dst, err := os.Create(tempFilePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to save the file: %v", err)})
		return
	}
	defer cleanTemp(tempFilePath, dst)

	_, err = io.Copy(dst, c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to save the file: %v", err)})
		return
	}

	guid := uuid.New().String()
	filePath, err := saveImageFromRequest(tempFilePath, guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Unable to save the file: %v", err)})
		return
	}

	uploadUrl, err := r2.GenerateUploadURL(guid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error generating R2 upload URL: %v", err)})
		return
	}

	err = uploadToUrl(uploadUrl, filePath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Error uploading to R2 upload URL: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"guid": guid})
}

func cleanTemp(tempFilePath string, dst *os.File) {
	dst.Close()
	os.Remove(tempFilePath)
}

func saveImageFromRequest(tempFilePath string, guid string) (string, error) {

	srcImage, err := imaging.Open(tempFilePath)
	if err != nil {
		return "", err
	}

	srcImage = imaging.Resize(srcImage, 400, 0, imaging.Lanczos)

	outputPath := filepath.Join(config.Configs.Cache_Dir, guid+".jpg")
	file, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	err = imaging.Encode(file, srcImage, imaging.JPEG, imaging.JPEGQuality(80))
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func uploadToUrl(uploadUrl *string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PUT", *uploadUrl, &fileBuffer)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "image/jpeg")
	req.ContentLength = int64(fileBuffer.Len())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func GetCachedImage(c *gin.Context) {
	objectKey := c.Query("key") // Assuming the object key is passed as a query parameter
	if objectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'key' query parameter"})
		return
	}

	imageData, found := cache.GetImage(objectKey)

	if !found {
		downloadUrl, err := generateDownloadURL(objectKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		resp, err := http.Get(downloadUrl)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to download image: %v", err)})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("unexpected status code: %d", resp.StatusCode)})
			return
		}

		imageData, err = io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read image data: %v", err)})
			return
		}

		if len(imageData) == 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "empty image file returned from R2"})
			return
		}

		cache.AddImage(objectKey, imageData)
	}

	c.Header("Content-Type", "image/jpeg")
	c.Header("Content-Length", fmt.Sprintf("%d", len(imageData)))
	c.Data(http.StatusOK, "image/jpeg", imageData)
	log.Printf("Wrote %d bytes to response", len(imageData))
}
