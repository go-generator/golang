package message

import (
	"time"

	. "../model"
	"github.com/common-go/search"
)

type TripSM struct {
	*search.SearchModel
	TripId    string         `json:"tripId"`
	StartTime *time.Time     `json:"startTime"`
	EndTime   *time.Time     `json:"endTime"`
	Locations []TripLocation `json:"locations"`
}
