package impl

import (
	"../../model"
	. "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type TripServiceImpl struct {
	*DefaultGenericService
	*DefaultSearchService
}

func NewTripServiceImpl(db *mongo.Database, searchBuilder SearchResultBuilder) *TripServiceImpl {
	var model model.Trip
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := NewMongoGenericSearchService(db, typeOfModel, "trips", searchBuilder, true, "")
	return &TripServiceImpl{genericService, searchService}
}
