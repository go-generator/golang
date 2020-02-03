package db_relationship

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	. "../cache/yaml/cache_cipher"
	. "../db_config"
	. "./constants"
	. "./database_models"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sqweek/dialog"
)

// 1-1 -> both fields are unique
// 1-n -> only one field is unique
// n-n -> both fields are not unique
// self reference will be in the same table with the same datatype

func DatabaseRelationships(dbConfig DatabaseConfig, conn *gorm.DB) ([]RelationshipTables, []string) {
	var dbConnect DatabaseConnection
	dbConnect.SetConnection(conn)
	rt := NewRelationshipTables(&dbConfig, &dbConnect)
	jt := ListAllJoinTablesWithCompositeKey(dbConfig.Database, dbConnect.GetConnection(), rt)
	for i := range rt {
		rt[i].Relationship = FindRelationShip(dbConfig.Database, &dbConnect, jt, &rt[i])
	}
	joinTable := ListAllJoinTablesWithCompositeKey(dbConfig.Database, dbConnect.GetConnection(), rt)
	return rt, joinTable
}

func JsonDescriptionGenerator(env, output string, conn *gorm.DB, dc *DatabaseConfig, rt []RelationshipTables) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var s SqlTablesData
	var files FilesDetails
	files.Env = []string{"search_model", "config", "controller", "service/impl", "route", "service"}
	var folderOutput Folders
	var entity []string
	if env == "" {
		env = "model"
	}
	tables := ListAllTableNames(conn, dc.Database)
	defer s.FreeResources() // Close connection before freeing resources
	files.Model = env
	s.InitSqlTablesData()
	path := "./" + output + "/"
	fileDirectory := path + output + ".json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
	}
	for _, v := range tables {
		var m ModelJSON
		m.Name = v
		s.InitSqlTable(dc.Database, v, conn)
		s.StandardizeFieldsName()
		for i, k := range s.SqlTable {
			var f FieldElements
			if s.ContainCompositeKey {
				f.Source = ToLower(s.GoFields[i])
			} else {
				if k.ColumnKey == "PRI" {
					f.Source = "_id"
				} else {
					f.Source = ToLower(s.GoFields[i])
				}
			}
			f.Type = s.TypeMap.TypeConvert[k.DataType]
			f.Name = s.GoFields[i]
			if k.ColumnKey == "PRI" {
				f.PrimaryKey = true
			} else {
				f.PrimaryKey = false
			}
			rl := GetRelationship(k.ColumnName, rt)
			if rl != nil {
				var foreign FieldElements
				if rl.Relationship == MTO && k.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					var relationship Relationship
					relationship.ReType = MTO
					relationship.Ref.Table = k.TableName
					foreign.Name = StandardizeName(rl.Table)
					foreign.Source = rl.Table
					foreign.Type = "*[]" + StandardizeName(rl.Table)
					foreign.ForeignKey = rl.Column
					relationship.Ref.RefCols = append(relationship.Ref.RefCols, rl.Column)
					m.Relationships = append(m.Relationships, relationship)
					m.Fields = append(m.Fields, foreign)
				}
				if rl.Relationship == MTO {
					f.ForeignKey = rl.ReferencedTable
				}
			}
			m.Fields = append(m.Fields, f)
		}
		files.Files = append(files.Files, m)
		s.ResetData() // Reuse Variable
	}
	for _, v := range files.Files {
		entity = append(entity, StandardizeName(v.Name))
	}
	files.Entity = entity
	folderOutput.ModelFile = append(folderOutput.ModelFile, files)
	data, err := json.MarshalIndent(&folderOutput, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fileDirectory, data, 0644) // Create and write files
	if err != nil {
		return err
	}
	return err
}

func JsonUI(env, filePath string, conn *gorm.DB, dc *DatabaseConfig, rt []RelationshipTables) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var s SqlTablesData
	var err error
	var files FilesDetails
	files.Env = []string{"search_model", "config", "controller", "service/impl", "route", "service"}
	var output Folders
	var entity []string
	tables := ListAllTableNames(conn, dc.Database)
	defer s.FreeResources() // Close connection before freeing resources
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	s.InitSqlTablesData()
	files.Model = env
	for _, v := range tables {
		var m ModelJSON
		m.Name = v
		s.InitSqlTable(dc.Database, v, conn)
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
			f.Type = s.TypeMap.TypeConvert[v.DataType]
			if v.ColumnKey == "PRI" {
				f.PrimaryKey = true
			} else {
				f.PrimaryKey = false
			}
			rl := GetRelationship(v.ColumnName, rt)
			if rl != nil {
				log.Println(rl)
				var relationship Relationship
				relationship.ReType = rl.Relationship
				var foreign FieldElements
				if rl.Relationship == MTO && v.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					relationship.Ref.Table = rl.Table
					foreign.Name = StandardizeName(rl.Table)
					foreign.Source = rl.Table
					foreign.Type = "*[]" + StandardizeName(rl.Table)
					foreign.ForeignKey = rl.Column
					relationship.Ref.RefCols = append(relationship.Ref.RefCols, rl.Column)
					m.Relationships = append(m.Relationships, relationship)
					m.Fields = append(m.Fields, foreign)
				}
				if rl.Relationship == MTO {
					f.ForeignKey = rl.ReferencedTable
				}
			}
			m.Fields = append(m.Fields, f)
		}
		files.Files = append(files.Files, m)
		s.ResetData() // Reuse Variable
	}
	for _, v := range files.Files {
		entity = append(entity, StandardizeName(v.Name))
	}
	files.Entity = entity
	output.ModelFile = append(output.ModelFile, files)
	data, err := json.MarshalIndent(&output, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, data, 0644) // Create and write files
	if err != nil {
		return err
	}
	return err
}

func InputValidation(app fyne.App, dc *DatabaseConfig, conn *gorm.DB) {
	if dc.Dialect == "" {
		ShowWindows(app, "Error", "Invalid Dialect")
		return
	}
	if dc.User == "" {
		ShowWindows(app, "Error", "Invalid User")
		return
	}
	if dc.Password == "" {
		ShowWindows(app, "Error", "Invalid Password")
		return
	}
	if dc.Host == "" {
		ShowWindows(app, "Error", "Invalid Host Address")
		return
	}
	if _, err := strconv.Atoi(strconv.Itoa(dc.Port)); err != nil {
		ShowWindows(app, "Error", "Invalid Port")
		return
	}
	if dc.Database == "" {
		ShowWindows(app, "Error", "Invalid Database Name")
		return
	}
	filename, errFile := dialog.File().Filter("JSON files", "json").Title("Save As").Save()
	if errFile != nil {
		ShowWindows(app, "Error", errFile.Error())
		return
	}
	tokens := strings.Split(filename, string(os.PathSeparator))
	rl, _ := DatabaseRelationships(*dc, conn)
	err := JsonUI(tokens[len(tokens)-2], filename+".json", conn, dc, rl)
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	ShowWindows(app, "Success", "Generated Database Json Description Successfully")
}

func InputUI(dc *DatabaseConfig, app fyne.App, cache, encryptField string) fyne.Window {
	var temp DatabaseConfig
	err := ReadCacheFile(cache, &temp, encryptField)
	if err != nil {
		log.Println(err)
	}
	window := app.NewWindow("Database Json Generator")
	window.Resize(fyne.Size{
		Width: 640,
	})
	//dialectEntry := widget.NewEntry()
	//dialectEntry.OnChanged = dc.SetDialect
	//dialectEntry.Text = dc.Dialect
	usernameEntry := widget.NewEntry()
	usernameEntry.OnChanged = dc.SetUsername
	usernameEntry.Text = dc.User
	passwordEntry := widget.NewEntry()
	passwordEntry.OnChanged = dc.SetPassword
	passwordEntry.Text = dc.Password
	passwordEntry.Password = true
	hostEntry := widget.NewEntry()
	hostEntry.OnChanged = dc.SetHost
	hostEntry.Text = dc.Host
	portEntry := widget.NewEntry()
	portEntry.OnChanged = dc.SetPort
	portEntry.Text = strconv.Itoa(dc.Port)
	databaseEntry := widget.NewEntry()
	databaseEntry.OnChanged = dc.SetDatabaseName
	databaseEntry.Text = dc.Database
	executeButton := widget.NewButton("Generate Database Json Description", func() {
		conn, err := dc.ConnectToSqlServer()
		if err != nil {
			log.Println(err)
		}
		defer func() {
			err = conn.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		InputValidation(app, dc, conn)
		if temp != *dc {
			err := WriteCacheFile(cache, dc, encryptField)
			if err != nil {
				log.Println(err)
			}
		}
	})
	providerEntry := widget.NewRadio([]string{"mysql", "postgres", "mssql", "sqlite3"}, func(s string) {
		dc.SetDialect(s)
	})
	window.SetContent(widget.NewVBox(
		//widget.NewLabel("Dialect:"),
		//dialectEntry,
		widget.NewLabel("Provider:"),
		providerEntry,
		widget.NewLabel("User:"),
		usernameEntry,
		widget.NewCheck("Show Password", func(b bool) {
			if b {
				passwordEntry.Password = false
			} else {
				passwordEntry.Password = true
			}
			passwordEntry.Refresh()
		}),
		widget.NewLabel("Password:"),
		passwordEntry,
		widget.NewLabel("Host:"),
		hostEntry,
		widget.NewLabel("Port:"),
		portEntry,
		widget.NewLabel("Database:"),
		databaseEntry,
		executeButton,
		widget.NewButton("Quit", func() {
			window.Close()
		}),
	))
	window.CenterOnScreen()
	return window
}
