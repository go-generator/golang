package model
import "github.com/go-generator/metadata"

type Output struct {
	ProjectName string          `json:"projectName"`
	RootPath    string          `json:"rootPath"`
	Files       []metadata.File `json:"files"`
	OutFile     []FileInfo
}
