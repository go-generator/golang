package controller

import (
	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"reflect"
)

type LocationRateController struct {
	*ViewController
	*SearchController
}

func NewLocationRateController(locationRateService LocationRateService, logService ActivityLogService) *LocationRateController {
	modelType := reflect.TypeOf(LocationRate{})
	searchModelType := reflect.TypeOf(LocationRateSM{})
	viewController := NewViewController(locationRateService, modelType, logService, "")
	searchController := NewSearchController(locationRateService, searchModelType, logService, false, "")
	return &LocationRateController{viewController, searchController}
}
