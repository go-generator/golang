package impl

import (
	"reflect"

	"../../model"
	. "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type EventServiceImpl struct {
	*DefaultGenericService
	*DefaultSearchService
}

func NewEventServiceImpl(db *mongo.Database, searchBuilder SearchResultBuilder) *EventServiceImpl {
	var model model.Event
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := NewMongoGenericSearchService(db, typeOfModel, "event", searchBuilder, true, "")
	return &EventServiceImpl{genericService, searchService}
}
