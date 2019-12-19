package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"reflect"

	. "../model"
	. "../search-model"
	. "../service"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
	"github.com/labstack/echo"
)

type EventController struct {
	*GenericController
	*SearchController
	EventService EventService
}

func NewEventController(eventService EventService, validator validator.Validator, logService ActivityLogService) *EventController {
	modelType := reflect.TypeOf(Event{})
	searchModelType := reflect.TypeOf(EventSM{})
	genericController := NewGenericController(eventService, modelType, nil, validator, logService, "")
	searchController := NewSearchController(eventService, searchModelType, logService, false, "")
	return &EventController{genericController, searchController, eventService}
}

func (c *EventController) GetEventByLocation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetEventByLocation", url)
		locationId := ctx.Param("locationId")
		list, err := c.EventService.GetEventByLocation(ctx.Request().Context(), locationId)
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
	}

}

func (c *EventController) GetEventByDate() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetEventByLocation", url)
		body := c.NewModel(ctx.Request().Body)
		list, err := c.EventService.GetEventByDate(ctx.Request().Context(), body["date"].(string))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *EventController) NewModel(body interface{}) (out map[string]interface{}) {
	if body != nil {
		switch s := body.(type) {
		case io.Reader:
			err := json.NewDecoder(s).Decode(&out)
			if err != nil {
				log.Println(err)
				return nil
			}
			return out
		}
	}
	return nil
}
