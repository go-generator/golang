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

	//TODO
	locationController := applicationContext.LocationController
	locationPath := "/location"
	e.GET(locationPath, locationController.GetAll())
	e.GET(locationPath+"/:id", locationController.GetById())
	e.POST(locationPath+"/search", locationController.Search())
	e.GET(locationPath+"/getByUrlId/:urlId", locationController.GetByUrlId())
	e.POST(locationPath+"/rateLocation", locationController.RateLocation())
	e.GET(locationPath+"/saveLocation/:userId/:locationId", locationController.SaveLocation())
	e.GET(locationPath+"/removeLocation/:userId/:locationId", locationController.RemoveLocation())
	e.GET(locationPath+"/getLocationsOfUser/:userId", locationController.GetLocationsOfUser())
	e.GET(locationPath+"/getLocationByTypeInRadius/:type/:radius", locationController.GetLocationByTypeInRadius())

	locationRateController := applicationContext.LocationRateController
	locationRatePath := "/locationRate"
	e.POST(locationRatePath+"/search", locationRateController.Search())

	tourController := applicationContext.TourController
	tourPath := "/tour"
	e.GET(tourPath, tourController.GetAll())
	e.GET(tourPath+"/:id", tourController.GetById())
	e.POST(tourPath+"/search", tourController.Search())

	bookingController := applicationContext.BookingController
	bookingPath := "/booking"
	e.GET(bookingPath, bookingController.GetAll())
	e.GET(bookingPath+"/:id", bookingController.GetById())
	e.POST(bookingPath+"/search", bookingController.Search())
	e.GET(bookingPath+"/cancel/:bookingId", bookingController.Cancel())
	e.POST(bookingPath+"/saveDraft", bookingController.SaveDraft())
	e.POST(bookingPath+"/getFreeLocationByBookable", bookingController.GetFreeLocationByBookable())
	e.POST(bookingPath+"/getFreeLocation", bookingController.GetLocationFreeInTime())
	e.POST(bookingPath+"/getFreeLocationByBookableList", bookingController.GetFreeLocationByBookableList())
	e.POST(bookingPath+"/submit", bookingController.Submit())

	bookableController := applicationContext.BookableController
	bookablePath := "/bookable"
	e.GET(bookablePath, bookableController.GetAll())
	e.GET(bookablePath+"/:id", bookableController.GetById())
	e.POST(bookablePath+"/search", bookableController.Search())

	eventController := applicationContext.EventController
	eventPath := "/event"
	e.GET(eventPath, eventController.GetAll())
	e.GET(eventPath+"/:id", eventController.GetById())
	e.POST(eventPath+"/search", eventController.Search())
	e.POST(eventPath, eventController.Insert())
	e.PUT(eventPath+"/:id", eventController.Update())
	e.PATCH(eventPath+"/partial/:id", eventController.Patch())
	e.GET(eventPath+"/getEventByLocation/:locationId", eventController.GetEventByLocation())
	e.POST(eventPath+"/getEventByDate", eventController.GetEventByDate())

	tripController := applicationContext.TripController
	tripPath := "/trip"
	e.GET(tripPath, tripController.GetAll())
	e.GET(tripPath+"/:id", tripController.GetById())
	e.POST(tripPath+"/search", tripController.Search())
	e.POST(tripPath, tripController.Insert())
	e.PUT(tripPath+"/:id", tripController.Update())

	return &TripalRoutes{e}, nil
}
