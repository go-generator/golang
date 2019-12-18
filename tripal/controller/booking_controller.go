package controller

import (
	. "../model"
	. "../search-model"
	. "../service"
	"encoding/json"
	"fmt"
	. "github.com/common-go/echo"
	"github.com/labstack/echo"
	"io"
	"log"
	"net/http"
	"reflect"
)

type BookingController struct {
	*ViewController
	*SearchController
	BookingService BookingService
}

func NewBookingController(bookingService BookingService, logService ActivityLogService) *BookingController {
	modelType := reflect.TypeOf(Booking{})
	searchModelType := reflect.TypeOf(BookingSM{})
	viewController := NewViewController(bookingService, modelType, logService, "")
	searchController := NewSearchController(bookingService, searchModelType, logService, false, "")
	return &BookingController{viewController, searchController, bookingService}
}

func (c *BookingController) Cancel() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to Cancel", url)
		bookingId := ctx.Param("bookingId")
		list, err := c.BookingService.Cancel(ctx.Request().Context(), bookingId)
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *BookingController) SaveDraft() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to SaveDraft", url)
		body := c.NewModel(ctx.Request().Body)
		fmt.Println("body", body)
		list, err := c.BookingService.SaveDraft(ctx.Request().Context(), body)
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *BookingController) GetFreeLocationByBookable() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetFreeLocationByBookable", url)
		body := c.NewModel1(ctx.Request().Body)
		list, err := c.BookingService.GetFreeLocationByBookable(ctx.Request().Context(), body["bookableId"].(string), body["date"].(string))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *BookingController) GetFreeLocationByBookableList() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetFreeLocationByBookableList", url)
		body := c.NewModel1(ctx.Request().Body)
		list, err := c.BookingService.GetFreeLocationByBookableList(ctx.Request().Context(), body["bookableIdList"].([]interface{}), body["date"].(string))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *BookingController) GetLocationFreeInTime() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to GetLocationFreeInTime", url)
		body := c.NewModel1(ctx.Request().Body)
		list, err := c.BookingService.GetLocationFreeInTime(ctx.Request().Context(), body["startDate"].(string), body["endDate"].(string))
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *BookingController) Submit() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to Submit", url)
		body := c.NewModel(ctx.Request().Body)
		list, err := c.BookingService.Submit(ctx.Request().Context(), body)
		if err != nil {
			ctx.Logger().Error(err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		} else {
			return ctx.JSON(http.StatusOK, list)
		}
		return nil
	}

}

func (c *BookingController) NewModel1(body interface{}) (out map[string]interface{}) {
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

func (c *BookingController) NewModel(body interface{}) (out Booking) {
	var result Booking
	//modelType := reflect.TypeOf(Booking{})
	//req := reflect.StatusNew(modelType).Interface()
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
