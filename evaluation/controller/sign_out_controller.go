package controller

import (
	//"../../core-echo/builder"
	. "github.com/common-go/echo"
	//"context"
	"fmt"
	"log"
	"time"
	//"context"
	//"time"
	"../service"
	"github.com/labstack/echo"
)

type SignOutController struct {
	signOutService 		service.SignOutService
	logService          ActivityLogService
}

func NewSignOutController(signOutService service.SignOutService, logService ActivityLogService) *SignOutController {
	return &SignOutController{signOutService, logService}
}

func (c *SignOutController) SignOut() echo.HandlerFunc{
	return func(ctx echo.Context) error {
		log.Println("Go to SignOut")

		fmt.Println("Token:", ctx.Get("token"))
		fmt.Println("ExpiresAt", ctx.Get("issuedAt").(time.Time))
		//list, err := c.UserRegistrationService.SignOut(ctx.Get("token").(string), ctx.Get("issuedAt").(time.Time))

		_, err := c.signOutService.SignOut(ctx.Get("token").(string), ctx.Get("issuedAt").(time.Time))
		if err != nil {
			//_ = Error(http.StatusInternalServerError, err, c.logService, ctx, c.Resource, "ImportArrayObject")
			return err
		}
		//_ = Succeed(http.StatusOK, list, c.logService, ctx, c.Resource, "ImportArrayObject")
		return nil
	}
}
