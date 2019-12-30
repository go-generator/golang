package main

import (
	"flag"
	"log"
	"strconv"

	. "./sql_data_model"
)

func RunWithCommandLine() {
	dialectPtr := flag.String("dialect", "", "input dialect")
	userPtr := flag.String("username", "", "input username")
	passPtr := flag.String("password", "", "input password")
	hostPtr := flag.String("host", "", "input host")
	portPtr := flag.String("port", "", "input port")
	dbNamePtr := flag.String("dbName", "", "input database name")
	packagePtr := flag.String("package", "", "input package name")
	outputPtr := flag.String("output", "", "input file output name")
	port, err := strconv.Atoi(*portPtr)
	if err != nil {
		log.Println(err)
	}
	flag.Parse()
	cre := DatabaseConfig{
		Dialect:  *dialectPtr,
		User:     *userPtr,
		Password: *passPtr,
		Host:     *hostPtr,
		Port:     port,
		Database: *dbNamePtr,
	}
	if *hostPtr == "" {
		cre.Host = "127.0.0.1:3306"
	}
	if *outputPtr == "" {
		*outputPtr = "database"
	}
	err = cre.JsonDescriptionGenerator(*packagePtr, *outputPtr)
	if err != nil {
		log.Println(err)
	}
}

func RunWithUI() {
	var cre DatabaseConfig
	cre.InputUI().ShowAndRun()
}

func main() {
	//RunWithCommandLine()
	RunWithUI()
}
