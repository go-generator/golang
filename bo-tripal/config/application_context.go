package config

import (
	"context"

	. "../builder"
	"../controller"
	"../service/impl"
	"github.com/common-go/mongo"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type ApplicationContext struct {
	LocationController     *controller.LocationController
	LocationRateController *controller.LocationRateController
	TourController         *controller.TourController
	BookingController      *controller.BookingController
	BookableController     *controller.BookableController
	EventController        *controller.EventController
}

func NewApplicationContext(mongoConfig mongo.MongoConfig) (*ApplicationContext, error) {
	ctx := context.Background()
	mongoDb, er1 := mongo.SetupMongo(ctx, mongoConfig)
	if er1 != nil {
		return nil, er1
	}

	mongoQueryBuilder := &mongo.DefaultQueryBuilder{}
	mongoSortBuilder := &mongo.DefaultSortBuilder{}

	mongoSearchResultBuilder := &mongo.DefaultSearchResultBuilder{
		Database:     mongoDb,
		QueryBuilder: mongoQueryBuilder,
		SortBuilder:  mongoSortBuilder,
	}
	locationSearchResultBuilder := &LocationSearchResultBuilder{
		Context:      ctx,
		Database:     mongoDb,
		QueryBuilder: mongoQueryBuilder,
		SortBuilder:  mongoSortBuilder,
	}
	// resultInfoBuilder := &builder.DefaultResultInfoBuilder{}

	//User activity Log mongo
	locationService := impl.NewLocationServiceImpl(mongoDb, locationSearchResultBuilder)
	locationController := controller.NewLocationController(locationService, nil, nil)

	locationRateService := impl.NewLocationRateServiceImpl(mongoDb, mongoSearchResultBuilder)
	locationRateController := controller.NewLocationRateController(locationRateService, nil)

	tourService := impl.NewTourServiceImpl(mongoDb, mongoSearchResultBuilder)
	tourController := controller.NewTourController(tourService, nil, nil)

	bookingService := impl.NewBookingServiceImpl(mongoDb, mongoSearchResultBuilder)
	bookingController := controller.NewBookingController(bookingService, nil)

	bookableService := impl.NewBookableServiceImpl(mongoDb, mongoSearchResultBuilder)
	bookableController := controller.NewBookableController(bookableService, nil, nil)

	eventService := impl.NewEventServiceImpl(mongoDb, mongoSearchResultBuilder)
	eventController := controller.NewEventController(eventService, nil, nil)

	return &ApplicationContext{locationController, locationRateController, tourController, bookingController, bookableController, eventController}, nil
}
