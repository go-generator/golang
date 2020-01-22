package cache

import (
	"log"
	"path/filepath"
)

const filePath = "./cache/cache.yaml"

var AbsPath = AbsPathCache()

func AbsPathCache() string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Println(err)
	}
	return absPath
}
