package main

import (
	c "./config"
	"./route"
	"github.com/common-go/config"
	"github.com/labstack/echo"
	"log"
	"os"
	"strconv"
)

func main() {
	parentPath := "tripal"
	resource := "resource"
	env := os.Getenv("ENV")
	var conf c.Root
	config.LoadConfig(parentPath, resource, env, &conf, "application")
	log.Println(" host ", conf)

	e := echo.New()
	_, er1 := route.NewTripalRoutes(e, conf.Mongo)
	if er1 != nil {
		panic(er1)
	}

	server := ""
	if conf.Server.Port > 0 {
		server = ":" + strconv.Itoa(conf.Server.Port)
	}
	e.Logger.Fatal(e.Start(server))
}
