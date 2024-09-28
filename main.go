package main

import (
	"log"
	"time"

	"r2-api-go/apihandlers"
	"r2-api-go/r2"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r2.Init()

	//router.HandleFunc("/generate-upload-url", apihandlers.GenerateUploadURL).Methods("GET")
	//router.HandleFunc("/generate-download-url", apihandlers.GenerateDownloadURL).Methods("GET")

	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true //TODO: config, do not allow all
	config.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	r.Use(cors.New(config))

	r.PUT("/upload-image", apihandlers.UploadImageGin)
	r.GET("/cached-image", apihandlers.GetCachedImageGin)

	port := ":8080"
	log.Printf("Server started at %s", port)
	log.Fatal(r.Run(port)) // Start the server on the specified port
}
