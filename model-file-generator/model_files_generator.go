package main

import (
	"flag"
	"log"
	"strconv"

	. "../database_config"
	. "../database_relationship"
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
	flag.Parse()
	port, err := strconv.Atoi(*portPtr)
	if err != nil {
		log.Println(err)
	}
	dbConfig := DatabaseConfig{
		Dialect:  *dialectPtr,
		User:     *userPtr,
		Password: *passPtr,
		Host:     *hostPtr,
		Port:     port,
		Database: *dbNamePtr,
	}
	if *hostPtr == "" {
		dbConfig.Host = "127.0.0.1:3306"
	}
	if *outputPtr == "" {
		*outputPtr = "database"
	}
	conn, err := dbConfig.ConnectToSqlServer()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	log.Println(DatabaseRelationships(dbConfig, conn))
	err = dbConfig.JsonDescriptionGenerator(*packagePtr, *outputPtr, conn)
	if err != nil {
		log.Println(err)
	}
}

func RunWithUI() {
	var cre DatabaseConfig
	cre.InputUI().ShowAndRun()
}

func main() {
	RunWithCommandLine()
	//RunWithUI()
}
