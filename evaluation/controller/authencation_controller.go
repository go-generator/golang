package controller

import (
	. "../../ldap-authentication"
	"encoding/json"
	"github.com/common-go/auth"
	. "github.com/common-go/echo"
	"github.com/labstack/echo"
	"log"
	"net/http"
)

type AuthenticationController struct {
	AuthenticationService AuthenticationService
	logService            ActivityLogService
}

func NewAuthenticationController(authenticationService AuthenticationService, logService ActivityLogService) *AuthenticationController {
	return &AuthenticationController{authenticationService, logService}
}

func (c *AuthenticationController) Authenticate() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		url := ctx.Path()
		log.Println("Go to Authentication", url)
		var user auth.AuthInfo
		er1 := json.NewDecoder(ctx.Request().Body).Decode(&user)
		if er1 != nil {
			return ctx.JSON(http.StatusBadRequest, "Error")
		}
		result, er2 := c.AuthenticationService.Authenticate(user)
		if er2 != nil {
			result.Status = auth.Fail
			log.Fatal(er2)
		}
		return ctx.JSON(http.StatusOK, result)
	}
}
