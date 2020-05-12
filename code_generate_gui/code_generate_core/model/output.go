package model
import iou "github.com/go-generator/io"

type Output struct {
	ProjectName string `json:"projectName"`
	RootPath    string `json:"rootPath"`
	Files []iou.File   `json:"files"`
	OutFile     []FileInfo
}
