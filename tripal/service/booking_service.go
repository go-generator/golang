package service

import (
	. "../model"
	"context"
	. "github.com/common-go/search"
	. "github.com/common-go/service"
)

type BookingService interface {
	ViewService
	SearchService
	Cancel(cxt context.Context, bookingId string) (bool, error)
	SaveDraft(cxt context.Context, booking Booking) (*Booking, error)
	GetFreeLocationByBookable(ctx context.Context, bookableId string, date string) ([]bool, error)
	Submit(ctx context.Context, objectBooking Booking) (*Booking, error)
	GetFreeLocationByBookableList(ctx context.Context, bookableIdList []interface{}, date string) ([]interface{}, error)
	GetLocationFreeInTime(ctx context.Context, startDate string, endDate string) ([]interface{}, error)
}
