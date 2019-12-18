package controller

import (
	"../handler"
	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
	"reflect"
)

type BookableController struct {
	*GenericController
	*SearchController
}

func NewBookableController(bookableService BookableService, validator validator.Validator, logService ActivityLogService) *BookableController {
	modelType := reflect.TypeOf(Bookable{})
	searchModelType := reflect.TypeOf(BookableSM{})
	idNames := bookableService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController := NewGenericController(bookableService, modelType, controlModelHandler, validator, logService, "")
	searchController := NewSearchController(bookableService, searchModelType, logService, false, "")
	return &BookableController{genericController, searchController}
}
