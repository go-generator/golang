package controller

import (
	"reflect"

	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
)

type EventController struct {
	*GenericController
	*SearchController
}

func NewEventController(eventService EventService, validator validator.Validator, logService ActivityLogService) *EventController {
	modelType := reflect.TypeOf(Event{})
	searchModelType := reflect.TypeOf(EventSM{})
	genericController := NewGenericController(eventService, modelType, nil, validator, logService, "")
	searchController := NewSearchController(eventService, searchModelType, logService, false, "")
	return &EventController{genericController, searchController}
}
