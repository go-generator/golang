package main

import (
	"flag"
	"log"

	. "../database_config"
	. "../database_relationship"
)

func RunWithCommandLine() {
	dialectPtr := flag.String("dialect", "", "input dialect")
	userPtr := flag.String("username", "", "input username")
	passPtr := flag.String("password", "", "input password")
	hostPtr := flag.String("host", "", "input host")
	portPtr := flag.Int("port", 0, "input port")
	dbNamePtr := flag.String("dbName", "", "input database name")
	packagePtr := flag.String("package", "", "input package name")
	outputPtr := flag.String("output", "", "input file output name")
	flag.Parse()
	dbConfig := DatabaseConfig{
		Dialect:  *dialectPtr,
		User:     *userPtr,
		Password: *passPtr,
		Host:     *hostPtr,
		Port:     *portPtr,
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
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	rl, _ := DatabaseRelationships(dbConfig, conn)
	err = JsonDescriptionGenerator(*packagePtr, *outputPtr, conn, &dbConfig, rl)
	if err != nil {
		log.Println(err)
	}
}

func RunWithUI() {
	var dbConfig DatabaseConfig
	InputUI(&dbConfig).ShowAndRun()
}

func main() {
	RunWithCommandLine()
	//RunWithUI()
}
