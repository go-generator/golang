package message

import (
	"time"

	. "../model"
	"github.com/common-go/search"
)

type BookingSM struct {
	*search.SearchModel
	BookingId        string         `json:"bookingId"`
	UserId           string         `json:"userId"`
	BookableId       string         `json:"bookableId"`
	Subject          string         `json:"subject"`
	Description      string         `json:"description"`
	StartBookingTime *time.Time     `json:"startBookingTime"`
	EndBookingTime   *time.Time     `json:"endBookingTime"`
	Status           BookableStatus `json:"status"`
}
