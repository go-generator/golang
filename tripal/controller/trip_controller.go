package controller

import (
	"reflect"

	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
)

type TripController struct {
	*GenericController
	*SearchController
}

func NewTripController(tripService TripService, validator validator.Validator, logService ActivityLogService) *TripController {
	modelType := reflect.TypeOf(Trip{})
	searchModelType := reflect.TypeOf(TripSM{})
	genericController := NewGenericController(tripService, modelType, nil, validator, logService, "")
	searchController := NewSearchController(tripService, searchModelType, logService, false, "")
	return &TripController{genericController, searchController}
}
