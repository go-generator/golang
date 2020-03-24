package impl

import (
	"reflect"

	"../../model"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type TourServiceImpl struct {
	*m.DefaultGenericService
	*m.DefaultSearchService
}

func NewTourServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *TourServiceImpl {
	var model model.Tour
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := m.NewMongoGenericSearchService(db, typeOfModel, "tours", searchBuilder, true, "")
	return &TourServiceImpl{genericService, searchService}
}
