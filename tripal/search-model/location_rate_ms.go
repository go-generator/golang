package message

import (
	"time"

	"github.com/common-go/search"
)

type LocationRateSM struct {
	*search.SearchModel
	RateId     string    `json:"rateId"`
	LocationId string    `json:"locationId"`
	UserId     string    `json:"userId"`
	Rate       int       `json:"rate"`
	RateTime   time.Time `json:"rateTime"`
	Review     string    `json:"review"`
}
