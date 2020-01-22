package map_type

import (
	"log"
	"path/filepath"
)

const filePath = "./map_type"

var DTypeAbsPath = DataTypeAbsPath()

func DataTypeAbsPath() string {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		log.Println(err)
	}
	return absPath
}
