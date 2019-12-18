package controller

import (
	"../handler"
	"../model"
	"../search-model"
	"../service"
	. "github.com/common-go/echo"
	//"context"
	//"log"
	//"net/http"
	"reflect"
	//"github.com/labstack/echo"

)

type CandidateEvaluationController struct {
	*GenericController
	*SearchController
}

func NewCandidateEvaluationController(candidateEvaluationService service.CandidateEvaluationService, logService ActivityLogService) *CandidateEvaluationController {
	var candidateEvaluationModel model.CandidateEvaluation
	modelType := reflect.TypeOf(candidateEvaluationModel)
	searchModelType := reflect.TypeOf(search_model.CandidateEvaluationSM{})
	idNames := candidateEvaluationService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController, searchController:= NewGenericSearchController(candidateEvaluationService, modelType, controlModelHandler, candidateEvaluationService, searchModelType,nil, logService, true, "")
	return &CandidateEvaluationController{GenericController: genericController, SearchController: searchController}
}


