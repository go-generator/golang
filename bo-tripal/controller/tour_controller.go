package controller

import (
	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
	"reflect"
)

type TourController struct {
	*GenericController
	*SearchController
}

func NewTourController(tourService TourService, validator validator.Validator, logService ActivityLogService) *TourController {
	modelType := reflect.TypeOf(Tour{})
	searchModelType := reflect.TypeOf(TourSM{})
	//idNames := tourService.GetIdNames()
	//controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController := NewGenericController(tourService, modelType, nil, validator, logService, "")
	searchController := NewSearchController(tourService, searchModelType, logService, false, "")
	return &TourController{genericController, searchController}
}
