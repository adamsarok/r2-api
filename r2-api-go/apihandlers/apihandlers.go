package apihandlers

import (
	"encoding/json"
	"net/http"
	"r2-api-go/r2-api-go/r2"
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

func GenerateDownloadURL(w http.ResponseWriter, r *http.Request) {
	objectKey := r.URL.Query().Get("key")
	if objectKey == "" {
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		return
	}

	downloadURL, err := r2.GenerateDownloadURL(objectKey)
	if err != nil {
		http.Error(w, "Failed to generate download URL "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson("downloadURL", *downloadURL, w)
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
	//generateUploadURL(w, r);

	// img, err := imaging.Open("input.jpg")
	// if err != nil {
	//     panic(err)
	// }

	// // Resize the image to width 800 preserving the aspect ratio
	// img = imaging.Resize(img, 800, 0, imaging.Lanczos)

	// // Save the resulting image as JPEG
	// err = imaging.Save(img, "output.jpg")
	// if err != nil {
	//     panic(err)
	// }

}
