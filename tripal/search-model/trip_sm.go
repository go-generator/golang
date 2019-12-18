package message

import (
	. "../model"
	"github.com/common-go/search"
	"time"
)

type TripSM struct {
	*search.SearchModel
	TripId    string         `json:"tripId"`
	StartTime *time.Time     `json:"startTime"`
	EndTime   *time.Time     `json:"endTime"`
	Locations []TripLocation `json:"locations"`
}
