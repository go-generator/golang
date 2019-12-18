package controller

import (
	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"reflect"
)

type BookableController struct {
	*ViewController
	*SearchController
}

func NewBookableController(bookableService BookableService, logService ActivityLogService) *BookableController {
	modelType := reflect.TypeOf(Bookable{})
	searchModelType := reflect.TypeOf(BookableSM{})
	viewController := NewViewController(bookableService, modelType, logService, "")
	searchController := NewSearchController(bookableService, searchModelType, logService, false, "")
	return &BookableController{viewController, searchController}
}
