package impl

import (
	"reflect"

	"../../model"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingServiceImpl struct {
	*m.DefaultViewService
	*m.DefaultSearchService
}

func NewBookingServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *BookingServiceImpl {
	var model model.Booking
	typeOfModel := reflect.TypeOf(model)
	viewService, searchService := m.NewMongoViewSearchService(typeOfModel, db, "booking", searchBuilder, false)
	return &BookingServiceImpl{viewService, searchService}
}
