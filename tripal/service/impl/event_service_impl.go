package impl

import (
	"../../model"
	"context"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
	"time"
)

type EventServiceImpl struct {
	*m.DefaultGenericService
	*m.DefaultSearchService
}

const (
	layoutDate string = "2006-01-02"
)

func NewEventServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *EventServiceImpl {
	var model model.Event
	typeOfModel := reflect.TypeOf(model)
	genericService, searchService := m.NewMongoGenericSearchService(db, typeOfModel, "event", searchBuilder, true, "")
	return &EventServiceImpl{genericService, searchService}
}

func (e *EventServiceImpl) GetEventByLocation(ctx context.Context, locationId string) ([]model.Event, error) {
	var query = bson.M{
		"locationId": locationId,
	}
	var model1 []model.Event
	//typeOfModel := reflect.TypeOf(model1)
	_, err := m.FindAndDecode(ctx, e.Collection, query, &model1)
	if err != nil {
		return nil, err
	} else {
		return model1, nil
	}
}

func (e *EventServiceImpl) GetEventByDate(ctx context.Context, date string) ([]model.Event, error) {
	//location := time.Now().Location()
	location, _ := time.LoadLocation("UTC")
	fromDate, err := time.ParseInLocation(layoutDate, date, location)
	if err != nil {
		return nil, err
	}
	dateEnd := time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 23, 59, 59, 0, location)
	var query = bson.M{
		"startTime": bson.M{
			"$gte": fromDate,
			"$lte": dateEnd,
		},
	}
	var model1 []model.Event
	// typeOfModel := reflect.TypeOf(model1)
	_, err = m.FindAndDecode(ctx, e.Collection, query, &model1)
	if err != nil {
		return nil, err
	} else {
		return model1, nil
	}
}
