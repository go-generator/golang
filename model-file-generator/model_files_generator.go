package main

import (
	"flag"
	"log"

	. "./sql_data_model"
)

func RunWithCommandLine() {
	userPtr := flag.String("username", "", "input username")
	passPtr := flag.String("password", "", "input password")
	dbNamePtr := flag.String("dbName", "", "input database name")
	hostPtr := flag.String("host", "", "input host")
	packagePtr := flag.String("package", "", "input package name")
	outputPtr := flag.String("output", "", "input file output name")
	flag.Parse()
	cre := DatabaseCredentials{
		Username:     *userPtr,
		Password:     *passPtr,
		DatabaseName: *dbNamePtr,
		Package:      *packagePtr,
		Output:       *outputPtr,
	}
	if *hostPtr == "" {
		cre.Host = "127.0.0.1:3306"
	}
	if *outputPtr == "" {
		*outputPtr = "database"
	}
	err := JsonDescriptionGenerator(cre)
	if err != nil {
		log.Fatal(err)
	}
}

func RunWithUI() {
	var cre DatabaseCredentials
	cre.InputUI().ShowAndRun()
}

func main() {
	//RunWithCommandLine()
	RunWithUI()
}
