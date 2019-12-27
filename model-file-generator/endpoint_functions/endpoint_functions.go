package endpoint_functions

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	. "../model_files_functions"
	. "../sql_data_model"
	"github.com/jinzhu/gorm"
)

func ModelGoFilesGenerator(s *SqlTablesData, conn *gorm.DB, tables []string, packageName string) {
	s.InitSqlTablesData()
	path := "./" + packageName + "/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			log.Fatal("Failed attempt to create directory, " + err.Error())
		}
	}
	for _, v := range tables {
		fileDirectory := "./" + packageName + "/" + v + ".go"
		s.InitSqlTable(v, conn)
		s.StandardizeFieldsName()
		s.CreateContent(packageName)
		err := ioutil.WriteFile(fileDirectory, []byte(s.WriteFile.String()), 0777) // Create and write files
		if err != nil {
			log.Fatal("Failed attempt to write model file," + err.Error())
		}
		s.ResetData() // Reuse Variable
	}
	log.Println("Model files are generated successfully")
}

func JsonDescriptionGenerator(cre DatabaseCredentials, packageName, output string) { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var s SqlTablesData
	var files FilesDetails
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
	files.Env = packageName
	s.InitSqlTablesData()
	path := "./" + packageName + "/"
	fileDirectory := path + output + ".json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			log.Fatal("Failed attempt to create directory, " + err.Error())
		}
	}
	for _, v := range tables {
		var m ModelJSON
		m.Name = v
		s.InitSqlTable(v, conn)
		s.StandardizeFieldsName()
		for i, v := range s.SqlTable {
			var f FieldElements
			if s.ContainCompositeKey {
				f.Source = ToLower(s.GoFields[i])
			} else {
				if v.ColumnKey == "PRI" {
					f.Source = "_id"
				} else {
					f.Source = ToLower(s.GoFields[i])
				}
			}
			f.Name = s.GoFields[i]
			f.Type = s.TypeMap[v.DataType]
			if v.ColumnKey == "PRI" {
				f.PrimaryKey = true
			} else {
				f.PrimaryKey = false
			}
			m.Fields = append(m.Fields, f)
		}
		files.Files = append(files.Files, m)
		s.ResetData() // Reuse Variable
	}
	data, err := json.MarshalIndent(&files, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(fileDirectory, data, 0644) // Create and write files
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Json files are generated successfully")
}
