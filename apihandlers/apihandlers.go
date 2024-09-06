package apihandlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"r2-api-go/cache"
	"r2-api-go/r2"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
)

func GenerateUploadURL(w http.ResponseWriter, r *http.Request) {
	objectKey := r.URL.Query().Get("key")
	if objectKey == "" {
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		return
	}

	uploadURL, err := r2.GenerateUploadURL(objectKey)
	if err != nil {
		http.Error(w, "Failed to generate upload URL "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson("uploadURL", *uploadURL, w)
}

func getObjectKey(r *http.Request) (string, error) {
	objectKey := r.URL.Query().Get("key")
	if objectKey == "" {
		return "", errors.New("missing 'key' query parameter")
	}
	return objectKey, nil
}

func GenerateDownloadURL(w http.ResponseWriter, r *http.Request) {
	objectKey, err := getObjectKey(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	downloadURL, err := generateDownloadURL(objectKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson("downloadURL", downloadURL, w)
}

func generateDownloadURL(objectKey string) (string, error) {
	downloadURL, err := r2.GenerateDownloadURL(objectKey)
	if err != nil {
		return "", errors.New("failed to generate download URL")
	}
	return *downloadURL, nil
}

func writeJson(param string, value string, w http.ResponseWriter) {
	response := map[string]string{
		param: value,
	}

	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to generate JSON response", http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}

func UploadImage(w http.ResponseWriter, r *http.Request) {
	tempFilename := r.URL.Query().Get("fileName")

	if tempFilename == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	fmt.Println("Uploaded Filename:", tempFilename)
	tempFilePath := filepath.Join("uploads", tempFilename)

	dst, err := os.Create(tempFilePath)
	if err != nil {
		http.Error(w, "Unable to create the file on disk", http.StatusInternalServerError)
		return
	}
	defer cleanTemp(tempFilePath, dst)

	_, err = io.Copy(dst, r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to save the file: %v", err), http.StatusInternalServerError)
		return
	}

	guid := uuid.New().String()
	filePath := saveImageFromRequest(w, tempFilePath, guid)

	uploadUrl, err := r2.GenerateUploadURL(guid)
	if err != nil {
		http.Error(w, "Error generating R2 upload URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	uploadToUrl(w, uploadUrl, filePath)

	writeJson("objectKey", guid, w)
}

func cleanTemp(tempFilePath string, dst *os.File) {
	dst.Close()
	os.Remove(tempFilePath)
}

func saveImageFromRequest(w http.ResponseWriter, tempFilePath string, guid string) string {

	srcImage, err := imaging.Open(tempFilePath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to open the iamge file: %v", err), http.StatusInternalServerError)
		return ""
	}

	srcImage = imaging.Resize(srcImage, 400, 0, imaging.Lanczos)

	outputPath := filepath.Join("uploads", guid+".jpg")
	file, err := os.Create(outputPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create file:  %v", err), http.StatusInternalServerError)
		return ""
	}
	defer file.Close()

	//err = imaging.Save(srcImage, outputPath)
	err = imaging.Encode(file, srcImage, imaging.JPEG, imaging.JPEGQuality(80))
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to convert the image file: %v", err), http.StatusInternalServerError)
		return ""
	}

	return outputPath
}

func uploadToUrl(w http.ResponseWriter, uploadUrl *string, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		http.Error(w, "Error opening file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		http.Error(w, "Error reading file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("PUT", *uploadUrl, &fileBuffer)
	if err != nil {
		http.Error(w, "Error creating request:"+err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Content-Type", "image/jpeg")
	req.ContentLength = int64(fileBuffer.Len())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error sending request:"+err.Error(), http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error sending request:", http.StatusInternalServerError)
		return
	}
}

func GetCachedImage(w http.ResponseWriter, r *http.Request) {
	objectKey, err := getObjectKey(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	imageData, found := cache.GetImage(objectKey)

	if !found {
		downloadUrl, err := generateDownloadURL(objectKey)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp, err := http.Get(downloadUrl)
		if err != nil {
			http.Error(w, "failed to download image: %v"+err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			http.Error(w, fmt.Sprintf("unexpected status code: %d", resp.StatusCode), http.StatusInternalServerError)
			return
		}

		imageData, err = io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to read image data: %v", err), http.StatusInternalServerError)
			return
		}

		if imageData == nil {
			http.Error(w, "empty image file returned from R2", http.StatusInternalServerError)
			return
		}

		cache.AddImage(objectKey, imageData)
	}

	//imageData, found = cache.GetImage(objectKey)

	w.Header().Set("Content-Type", "image/jpeg") // Assuming JPEG for this example
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(imageData)))
	w.WriteHeader(http.StatusOK)

	bytesWritten, err := w.Write(imageData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error writing image data to response: %v", err), http.StatusInternalServerError)
		return
	}
	log.Printf("Wrote %d bytes to response", bytesWritten)
}
