package main

import (
	"log"
	"net/http"

	"r2-api-go/apihandlers"
	"r2-api-go/r2"

	"github.com/gorilla/mux"
)

func main() {
	r2.Init()

	router := mux.NewRouter()

	router.HandleFunc("/generate-upload-url", apihandlers.GenerateUploadURL).Methods("GET")
	router.HandleFunc("/generate-download-url", apihandlers.GenerateDownloadURL).Methods("GET")
	router.HandleFunc("/upload-image", apihandlers.UploadImage).Methods("PUT")

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
