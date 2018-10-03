package main

import (
	"log"
	"net/http"
	"github.com/julienschmidt/httprouter"
	. "./context"
	. "./controller"
)

var appContext = AppContext{}
var merchantController = MerchantController{}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	merchantController = appContext.GetMerchantController()
}

// Define HTTP request routes
func main() {
	router := httprouter.New()
	router.GET("/merchants", merchantController.GetAll)
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
