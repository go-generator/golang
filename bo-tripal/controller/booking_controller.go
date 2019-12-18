package controller

import (
	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"reflect"
)

type BookingController struct {
	*ViewController
	*SearchController
}

func NewBookingController(bookingService BookingService, logService ActivityLogService) *BookingController {
	modelType := reflect.TypeOf(Booking{})
	searchModelType := reflect.TypeOf(BookingSM{})
	viewController := NewViewController(bookingService, modelType, logService, "")
	searchController := NewSearchController(bookingService, searchModelType, logService, false, "")
	return &BookingController{viewController, searchController}
}
