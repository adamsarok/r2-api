package cache

import (
	lru "github.com/hashicorp/golang-lru"
)

var ( //package level variables are instantiated before anything
	imageCache = newImageCache(512)
)

// func GetImageCache() *ImageCache {
// 	return imageCache
// }

type ImageCache struct {
	cache *lru.Cache
}

func newImageCache(size int) *ImageCache {
	cache, _ := lru.New(size)
	return &ImageCache{cache: cache}
}

func AddImage(key string, imageData []byte) {
	imageCache.cache.Add(key, imageData)
}

func GetImage(key string) ([]byte, bool) {
	image, found := imageCache.cache.Get(key)
	if found {
		return image.([]byte), true
	}
	return nil, false
}
