package controller

import (
	"../handler"
	"../model"
	"../search-model"
	"../service"
	. "github.com/common-go/echo"
	"reflect"
)

type BatchController struct {
	*GenericController
	*SearchController
}


func NewBatchController(batchService service.BatchService, logService ActivityLogService) *BatchController {
	var batchModel model.Batch
	modelType := reflect.TypeOf(batchModel)
	searchModelType := reflect.TypeOf(search_model.BatchSM{})
	idNames := batchService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController, searchController:= NewGenericSearchController(batchService, modelType, controlModelHandler, batchService, searchModelType,nil, logService, true, "")
	return &BatchController{GenericController: genericController, SearchController: searchController}
}
