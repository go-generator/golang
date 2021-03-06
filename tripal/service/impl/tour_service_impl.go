package impl

import (
	"reflect"

	"../../model"
	. "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourServiceImpl struct {
	*DefaultViewService
	*DefaultSearchService
}

func NewTourServiceImpl(db *mongo.Database, searchBuilder SearchResultBuilder) *TourServiceImpl {
	var model model.Tour
	typeOfModel := reflect.TypeOf(model)
	viewService, searchService := NewMongoViewSearchService(typeOfModel, db, "tours", searchBuilder, true)
	return &TourServiceImpl{viewService, searchService}
}
