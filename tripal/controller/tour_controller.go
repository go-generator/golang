package controller

import (
	"reflect"

	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
)

type TourController struct {
	*ViewController
	*SearchController
}

func NewTourController(tourService TourService, logService ActivityLogService) *TourController {
	modelType := reflect.TypeOf(Tour{})
	searchModelType := reflect.TypeOf(TourSM{})
	viewController := NewViewController(tourService, modelType, logService, "")
	searchController := NewSearchController(tourService, searchModelType, logService, false, "")
	return &TourController{viewController, searchController}
}
