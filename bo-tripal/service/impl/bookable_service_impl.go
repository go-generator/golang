package impl

import (
	"reflect"

	"../../model"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookableServiceImpl struct {
	*m.DefaultGenericService
	*m.DefaultSearchService
}

func NewBookableServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *BookableServiceImpl {
	var model model.Bookable
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := m.NewMongoGenericSearchService(db, typeOfModel, "bookable", searchBuilder, false, "")
	return &BookableServiceImpl{genericService, searchService}
}
