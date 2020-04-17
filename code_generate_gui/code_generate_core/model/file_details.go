package model

import . "github.com/go-generator/metadata"

type FilesDetails struct {
	Model string  `json:"model"`
	Files []Model `json:"files"`
}
