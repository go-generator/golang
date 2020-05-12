package model
import "github.com/go-generator/metadata"

type Input struct {
	Folders []metadata.Project `json:"folders"`
}
