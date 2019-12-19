package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	c "./config"
	. "./route"
	"github.com/common-go/config"
	server "github.com/common-go/echo"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func StartServerTLS(handler http.Handler, cert string, key string, conf server.ServerConfig) {
	serverPort := ""
	if conf.Port > 0 {
		serverPort = ":" + strconv.Itoa(conf.Port)
	}
	if len(conf.Version) > 0 {
		if conf.Port > 0 {
			log.Println("Start service: " + conf.Name + " at port " + strconv.Itoa(conf.Port) + " with version " + conf.Version)
		} else {
			log.Println("Start service: " + conf.Name + " with version " + conf.Version)
		}
	} else {
		if conf.Port > 0 {
			log.Println("Start service: " + conf.Name + " at port " + strconv.Itoa(conf.Port))
		} else {
			log.Println("Start service: " + conf.Name)
		}
	}
	if err := http.ListenAndServeTLS(serverPort, cert, key, handler); err != nil {
		panic(err)
	}
}

func StartEchoTLS(e *echo.Echo, cert []byte, key []byte, conf server.ServerConfig) {
	serverPort := ""
	if conf.Port > 0 {
		serverPort = ":" + strconv.Itoa(conf.Port)
	}
	if len(conf.Version) > 0 {
		if conf.Port > 0 {
			log.Println("Start service: " + conf.Name + " at port " + strconv.Itoa(conf.Port) + " with version " + conf.Version)
		} else {
			log.Println("Start service: " + conf.Name + " with version " + conf.Version)
		}
	} else {
		if conf.Port > 0 {
			log.Println("Start service: " + conf.Name + " at port " + strconv.Itoa(conf.Port))
		} else {
			log.Println("Start service: " + conf.Name)
		}
	}
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"authorization", "Content-Type"},
		AllowCredentials: true,
		AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Logger.Fatal(e.StartTLS(serverPort, cert, key))
}

func main() {
	parentPath := "evaluation"
	resource := "resource"
	env := os.Getenv("ENV")
	var conf c.Root
	config.LoadConfig(parentPath, resource, env, &conf, "application")
	//cert, err := config.LoadFile(parentPath, resource, "cert.pem") // Read certificate file in resource
	//if err != nil {
	//	panic(err)
	//}
	//key, err := config.LoadFile(parentPath, resource, "key.pem") // Read key file in resource
	//if err != nil {
	//	panic(err)
	//}
	//resourceMap := server.LoadMap(parentPath, resource, env, "resource_map")
	log.Println(" host ", conf)
	e := echo.New()
	route, err := NewEvaluationRoutes(e, conf.Mongo, conf.Ldap, conf.Token)
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
