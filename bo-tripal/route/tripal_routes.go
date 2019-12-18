package route

import (
	"../config"
	"github.com/common-go/mongo"
	"github.com/labstack/echo"
)

type TripalRoutes struct {
	Echo *echo.Echo
}

func NewTripalRoutes(e *echo.Echo, mongoConfig mongo.MongoConfig) (*TripalRoutes, error) {
	applicationContext, err := config.NewApplicationContext(mongoConfig)
	if err != nil {
		return nil, err
	}
	/*
		router.POST(bankPath, bankController.Insert)
		router.PUT(bankPath+"/:id", bankController.Update)
		router.PATCH(bankPath+"/partial"+"/:id", bankController.UpdatePartial)
		router.DELETE(bankPath+"/:id", bankController.Delete)
		router.DELETE(bankPath+"/:id", bankController.Delete)
	*/

	locationController := applicationContext.LocationController
	locationPath := "/location"
	e.GET(locationPath, locationController.GetAll())
	e.GET(locationPath+"/:id", locationController.GetById())
	e.POST(locationPath+"/search", locationController.Search())
	e.POST(locationPath, locationController.Insert())
	e.PUT(locationPath+"/:id", locationController.Update())

	locationRateController := applicationContext.LocationRateController
	locationRatePath := "/locationRate"
	e.GET(locationRatePath, locationRateController.GetAll())
	e.GET(locationRatePath+"/:id", locationRateController.GetById())
	e.POST(locationRatePath+"/search", locationRateController.Search())

	tourController := applicationContext.TourController
	tourPath := "/tour"
	e.GET(tourPath, tourController.GetAll())
	e.GET(tourPath+"/:id", tourController.GetById())
	e.POST(tourPath+"/search", tourController.Search())
	e.POST(tourPath, tourController.Insert())
	e.PUT(tourPath+"/:id", tourController.Update())

	bookingController := applicationContext.BookingController
	bookingPath := "/booking"
	e.GET(bookingPath, bookingController.GetAll())
	e.GET(bookingPath+"/:id", bookingController.GetById())
	e.POST(bookingPath+"/search", bookingController.Search())

	bookableController := applicationContext.BookableController
	bookablePath := "/bookable"
	e.GET(bookablePath, bookableController.GetAll())
	e.GET(bookablePath+"/:id", bookableController.GetById())
	e.POST(bookablePath+"/search", bookableController.Search())
	e.POST(bookablePath, bookableController.Insert())
	e.PUT(bookablePath+"/:id", bookableController.Update())

	eventController := applicationContext.EventController
	eventPath := "/event"
	e.GET(eventPath, eventController.GetAll())
	e.GET(eventPath+"/:id", eventController.GetById())
	e.POST(eventPath+"/search", eventController.Search())
	e.POST(eventPath, eventController.Insert())
	e.PUT(eventPath+"/:id", eventController.Update())

	return &TripalRoutes{e}, nil
}
