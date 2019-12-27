package main

import (
	"flag"

	. "./endpoint_functions"
	. "./sql_data_model"
)

func main() {
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
	}
	if *hostPtr == "" {
		cre.Host = "127.0.0.1:3306"
	}
	if *outputPtr == "" {
		*outputPtr = "database"
	}
	JsonDescriptionGenerator(cre, *packagePtr, *outputPtr) // Generate database description json file
}
