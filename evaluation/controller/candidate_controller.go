package controller

import (
	"../handler"
	"../model"
	"../search-model"
	"../service"
	"context"
	. "github.com/common-go/echo"
	"github.com/labstack/echo"
	//"encoding/json"
	//"log"
	//"net/http"
	"reflect"
	//"encoding/json"

)

type CandidateController struct {
	*GenericController
	*SearchController
	logService     ActivityLogService
	CandidateService service.CandidateService
}

func NewCandidateController(candidateService service.CandidateService, logService ActivityLogService) *CandidateController {
	var candidateModel model.Candidate
	modelType := reflect.TypeOf(candidateModel)
	searchModelType := reflect.TypeOf(search_model.CandidateSM{})
	idNames := candidateService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController, searchController:= NewGenericSearchController(candidateService, modelType, controlModelHandler, candidateService, searchModelType,nil, logService, true, "")
	return &CandidateController{GenericController: genericController, SearchController: searchController}
}

func (c *CandidateController) ImportArrayObject() echo.HandlerFunc{
	return func(ctx echo.Context) error {
		var candidate []model.Candidate
		//var r *http.Request
		err :=  ctx.Bind(&candidate)
		//candidate, _ = url.ParseQuery(url_candidate)
		var context  context.Context
		_, err1 := c.CandidateService.ImportArrayObject(context, candidate)
		if err != nil {
			//_ = Error(http.StatusInternalServerError, err, c.logService, ctx, c.Resource, "ImportArrayObject")
			return err
		}
		//_ = Succeed(http.StatusOK, list, c.logService, ctx, c.Resource, "ImportArrayObject")
		return err1
	}
}

//func (c *CandidateController) PatchMark(ctx echo.Context) error{
//	url := ctx.Request().URL.String()
//	log.Println("Go to PatchMark", url)
//
//	var id string
//	id = ctx.Param("_id")
//	var mark float64
//	mark = ctx.Param("mark")
//
//	err := c.CandidateService.PatchMark(ctx, id, mark)
//	if err != nil {
//		//_ = Error(http.StatusInternalServerError, err, c.logService, ctx, c.Resource, "PatchMark")
//		return err
//	}
//	//_ = Succeed(http.StatusOK, "", c.logService, ctx, c.Resource, "PatchMark")
//	return nil
//}





