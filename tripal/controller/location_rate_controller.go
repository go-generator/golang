package controller

import (
	"reflect"

	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
)

type LocationRateController struct {
	*SearchController
}

func NewLocationRateController(locationRateService LocationRateService, logService ActivityLogService) *LocationRateController {
	searchModelType := reflect.TypeOf(LocationRateSM{})
	searchController := NewSearchController(locationRateService, searchModelType, logService, false, "")
	return &LocationRateController{searchController}
}
