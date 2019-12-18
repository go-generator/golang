package message

import (
	. "../model"
	"github.com/common-go/search"
)

type BookableSM struct {
	*search.SearchModel
	BookableId          string       `json:"bookableId"`
	LocationId          string       `json:"locationId"`
	BookableType        BookableType `json:"bookableType"`
	BookableName        string       `json:"bookableName"`
	BookableDescription string       `json:"bookableDescription"`
	BookableCapacity    int          `json:"bookableCapacity"`
	Image               string       `json:"image"`
}
