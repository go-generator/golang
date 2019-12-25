package main

import (
	"flag"
	"log"

	. "golang/model-file-generator/model_files_functions"
	. "golang/model-file-generator/sql_data_model"
)

func main() {
	var s SqlTablesData
	userPtr := flag.String("username", "", "input username")
	passPtr := flag.String("password", "", "input password")
	dbNamePtr := flag.String("dbName", "", "input database name")
	hostPtr := flag.String("host", "", "input host")
	flag.Parse()
	packageName := "model"
	cre := DatabaseCredentials{
		Username:     *userPtr,   //"test",
		Password:     *passPtr,   //"Doraemon1096~", //"127.0.0.1:3306",
		DatabaseName: *dbNamePtr, //"odd",
	}
	if *hostPtr == "" {
		cre.Host = "127.0.0.1:3306"
	}
	conn, err := cre.ConnectToSqlServer()
	log.Println("Connecting to sql server...")
	if err != nil {
		log.Fatal("Failed attempt to connect to server, " + err.Error())
	}
	log.Println("Connection to sql server is established successfully")
	tables := ListAllTableNames(conn, cre.DatabaseName)
	defer s.FreeResources() // Close connection before freeing resources
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal("Failed attempt to close the connection, " + err.Error())
		}
	}()
	ModelFilesGenerator(&s, conn, tables, packageName)
}
