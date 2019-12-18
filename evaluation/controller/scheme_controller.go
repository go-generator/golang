package controller


import (
	"../handler"
	"../model"
	"../search-model"
	"../service"
	. "github.com/common-go/echo"
	"reflect"
)

type SchemeController struct {
	*GenericController
	*SearchController
}


func NewSchemeController(schemeService service.SchemeService, logService ActivityLogService) *SchemeController {
	var schemeModel model.Scheme
	modelType := reflect.TypeOf(schemeModel)
	searchModelType := reflect.TypeOf(search_model.SchemeSM{})
	idNames := schemeService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController, searchController:= NewGenericSearchController(schemeService, modelType, controlModelHandler, schemeService, searchModelType,nil, logService, true, "")
	return &SchemeController{GenericController: genericController, SearchController: searchController}
}

