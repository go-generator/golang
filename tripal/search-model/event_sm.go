package message

import (
	"time"

	"github.com/common-go/search"
)

type EventSM struct {
	*search.SearchModel
	EventId    string     `json:"eventId"`
	EventName  string     `json:"eventName"`
	StartTime  *time.Time `json:"startTime"`
	EndTime    *time.Time `json:"endTime"`
	LocationId string     `json:"locationId"`
	Lat        float64    `json:"lat"`
	Long       float64    `json:"long"`
}
