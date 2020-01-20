package cache_file

import (
	"log"
	"path/filepath"
)

const filePath = "./cache_file/cache.yaml"

var AbsPath = CacheAbsPath()

func CacheAbsPath() string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Println(err)
	}
	return absPath
}
