package impl

import (
	"reflect"

	"../../model"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookableServiceImpl struct {
	*m.DefaultViewService
	*m.DefaultSearchService
}

func NewBookableServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *BookableServiceImpl {
	var model model.Bookable
	typeOfModel := reflect.TypeOf(model)
	viewService, searchService := m.NewMongoViewSearchService(typeOfModel, db, "bookable", searchBuilder, false)
	return &BookableServiceImpl{viewService, searchService}
}
