package main

import (
	"log"
	"os"
	"strconv"

	c "./config"
	"./route"
	"github.com/common-go/config"
	"github.com/labstack/echo"
)

func main() {
	parentPath := "tripal"
	resource := "resource"
	env := os.Getenv("ENV")
	var conf c.Root
	config.LoadConfig(parentPath, resource, env, &conf, "application")
	cert, err := config.LoadFile(parentPath, resource, "cert.pem")
	if err != nil {
		panic(err)
	}
	key, err := config.LoadFile(parentPath, resource, "key.pem")
	if err != nil {
		panic(err)
	}
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
	if conf.Server.Secure {
		e.Logger.Fatal(e.StartTLS(server, cert, key))
	} else {
		e.Logger.Fatal(e.Start(server))
	}
}
