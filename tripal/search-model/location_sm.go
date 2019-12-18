package message

import (
	"github.com/common-go/search"
)

type LocationSM struct {
	*search.SearchModel
	LocationId   string `json:"locationId"`
	Type         string `json:"type"`
	LocationName string `json:"locationName`
}
