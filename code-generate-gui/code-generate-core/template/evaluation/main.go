package main

import (
	c "evaluation/config"
	. "evaluation/route"
	"fmt"
	"github.com/common-go/config"
	"github.com/common-go/echo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"os"
)

func main() {
	parentPath := "evaluation"
	resource := "resource"
	env := os.Getenv("ENV")
	var conf c.Root
	config.LoadConfig(parentPath, resource, env, &conf, "application")
	//resourceMap := server.LoadMap(parentPath, resource, env, "resource_map")
	log.Println(" host ", conf)
	e := echo.New()
	route, err := NewEvaluationRoutes(e, conf.Mongo)
	if err != nil {
		panic(fmt.Errorf("create route failed"))
	}
	//evaRoutes, er1 := route.NewEvaluationRoutes(e, conf.Database, conf.Mongo, conf.Ldap, "secrettma", 86400000, resourceMap)
	//if er1 != nil {
	//	panic(er1)
	//}
	route.Router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete, http.MethodPatch},
	}))
	server.StartServer(route.Router, conf.Server)
}
