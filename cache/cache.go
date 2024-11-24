package cache

import (
	"log"
	"os"
	"path/filepath"
	"r2-api-go/config"

	lru "github.com/hashicorp/golang-lru"
)

var ( //package level variables are instantiated before anything
	inMemoryCache = newImageCache(512)
)

// func GetImageCache() *ImageCache {
// 	return imageCache
// }

type ImageCache struct {
	cache *lru.Cache
}

func getFilePath(key string) string {
	return filepath.Join(config.Configs.Cache_Dir, key+".jpg")
}

func newImageCache(size int) *ImageCache {
	cache, _ := lru.New(size)
	return &ImageCache{cache: cache}
}

func AddImage(key string, imageData []byte) {
	inMemoryCache.cache.Add(key, imageData)
	AddToDiskCache(key, imageData)
}

func GetImage(key string) ([]byte, bool) {
	image, found := inMemoryCache.cache.Get(key)
	if found {
		return image.([]byte), true
	}

	cachedPath := getFilePath(key)
	if !fileExists(cachedPath) {
		return nil, false
	}

	image, err := os.ReadFile(getFilePath(key))
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	return image.([]byte), true
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func AddToDiskCache(key string, imageData []byte) error {
	file, err := os.Create(filepath.Join(config.Configs.Cache_Dir, key+".jpg"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = file.Write(imageData)
	if err != nil {
		panic(err)
	}

	return nil
}
