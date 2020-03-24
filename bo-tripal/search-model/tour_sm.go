package message

import (
	"time"

	"github.com/common-go/search"
)

type TourSM struct {
	*search.SearchModel
	TourId    string     `json:"tourId"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	Locations []string   `json:"locations"`
}
