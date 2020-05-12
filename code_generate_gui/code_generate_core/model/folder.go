package model

import . "github.com/go-generator/metadata"

type Folder struct {
	Array []string	`json:"array"`
	Entity []string `json:"entity"`
	RawEnv []string `json:"env"`
	Model  string   `json:"model"`
	Files  []Model  `json:"files"`
}
