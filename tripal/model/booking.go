package model

import "time"

type BookableStatus string

const (
	News      BookableStatus = "N"
	Submitted BookableStatus = "S"
	Approved  BookableStatus = "A"
	Rejected  BookableStatus = "R"
	Cancelled BookableStatus = "C"
)

type Booking struct {
	BookingId        string         `json:"bookingId" bson:"_id"`
	UserId           string         `json:"userId" bson:"userId"`
	BookableId       string         `json:"bookableId" bson:"bookableId"`
	Subject          string         `json:"subject" bson:"subject"`
	Description      string         `json:"description" bson:"description"`
	StartBookingTime *time.Time     `json:"startBookingTime" bson:"startBookingTime"`
	EndBookingTime   *time.Time     `json:"endBookingTime" bson:"endBookingTime"`
	Status           BookableStatus `json:"status" bson:"status"`
}
