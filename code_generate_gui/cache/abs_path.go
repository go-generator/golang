package cache

import (
	"path/filepath"

	"golang/code_generate_gui/working_directory"
)

var AbsPath = AbsPathCache()

func AbsPathCache() string {
	filePath := []string{working_directory.GetWorkingDirectory(), "cache", "cache.yaml"}
	absPath := filepath.Join(filePath...)
	return absPath
}
