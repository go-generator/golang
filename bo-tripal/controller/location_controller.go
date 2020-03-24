package controller

import (
	"reflect"

	"../handler"
	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
)

type LocationController struct {
	*GenericController
	*SearchController
}

func NewLocationController(locationService LocationService, validator validator.Validator, logService ActivityLogService) *LocationController {
	modelType := reflect.TypeOf(Location{})
	searchModelType := reflect.TypeOf(LocationSM{})
	idNames := locationService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController := NewGenericController(locationService, modelType, controlModelHandler, validator, logService, "")
	searchController := NewSearchController(locationService, searchModelType, logService, false, "")
	return &LocationController{genericController, searchController}
}
