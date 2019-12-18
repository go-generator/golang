package message

import (
	"github.com/common-go/search"
	"time"
)

type TourSM struct {
	*search.SearchModel
	TourId    string     `json:"tourId"`
	StartTime *time.Time `json:"startTime"`
	EndTime   *time.Time `json:"endTime"`
	Locations []string   `json:"locations"`
}
