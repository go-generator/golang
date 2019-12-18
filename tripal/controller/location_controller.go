package controller

import (
	"../handler"
	. "../model"
	. "../search-model"
	. "../service"
	"encoding/json"
	. "github.com/common-go/echo"
	"github.com/common-go/validator"
	"github.com/labstack/echo"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
)

type LocationController struct {
	*GenericController
	*SearchController
	LocationService LocationService
}

func NewLocationController(locationService LocationService, validator validator.Validator, logService ActivityLogService) *LocationController {
	modelType := reflect.TypeOf(Location{})
	searchModelType := reflect.TypeOf(LocationSM{})
	idNames := locationService.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController := NewGenericController(locationService, modelType, controlModelHandler, validator, logService, "")
	searchController := NewSearchController(locationService, searchModelType, logService, false, "")
	return &LocationController{genericController, searchController, locationService}
}

func (c *LocationController) GetByUrlId() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetByUrlId", url)
		list, err := c.LocationService.GetByUrlId(ctx.Request().Context(), ctx.Param("urlId"))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *LocationController) RateLocation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to RateLocation", url)
		body := c.NewModel(ctx.Request().Body)
		list, err := c.LocationService.RateLocation(ctx.Request().Context(), body)
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *LocationController) SaveLocation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to SaveLocation", url)
		list, err := c.LocationService.SaveLocation(ctx.Request().Context(), ctx.Param("userId"), ctx.Param("locationId"))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}
}

func (c *LocationController) RemoveLocation() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to RemoveLocation", url)
		list, err := c.LocationService.RemoveLocation(ctx.Request().Context(), ctx.Param("userId"), ctx.Param("locationId"))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}
}

func (c *LocationController) GetLocationsOfUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetLocationsOfUser", url)
		list, err := c.LocationService.GetLocationsOfUser(ctx.Request().Context(), ctx.Param("userId"))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}
}

func (c *LocationController) GetLocationByTypeInRadius() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetLocationByTypeInRadius", url)
		radius, er1 := strconv.Atoi(ctx.Param("radius"))
		if er1 != nil {
			radius = 0
		}
		list, err := c.LocationService.GetLocationByTypeInRadius(ctx.Request().Context(), ctx.Param("type"), radius)
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}
}

func (c *LocationController) NewModel(body interface{}) (out LocationRate) {
	var result LocationRate
	//modelType := reflect.TypeOf(Booking{})
	//req := reflect.New(modelType).Interface()
	if body != nil {
		switch s := body.(type) {
		case io.Reader:
			err := json.NewDecoder(s).Decode(&result)
			if err != nil {
				log.Println(err)
				return result
			}
			return result
		}
	}
	return result
}
