package impl

import (
	"reflect"

	"../../model"
	. "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationRateServiceImpl struct {
	*DefaultViewService
	*DefaultSearchService
}

func NewLocationRateServiceImpl(db *mongo.Database, searchBuilder SearchResultBuilder) *LocationRateServiceImpl {
	var model model.LocationRate
	typeOfModel := reflect.TypeOf(model)
	viewService, searchService := NewMongoViewSearchService(typeOfModel, db, "locationRate", searchBuilder, true)
	return &LocationRateServiceImpl{viewService, searchService}
}
