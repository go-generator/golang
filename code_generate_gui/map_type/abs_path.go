package map_type

import (
	"path/filepath"

	"golang/code_generate_gui/working_directory"
)

var DTypeAbsPath = DataTypeAbsPath()

func DataTypeAbsPath() string {
	filePath := []string{working_directory.GetWorkingDirectory(), "map_type"}
	absPath := filepath.Join(filePath...)
	return absPath
}
