package database_relationship

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	. "../database_config"
	. "./constants"
	. "./database_models"
	"fyne.io/fyne"
	fApp "fyne.io/fyne/app"
	"fyne.io/fyne/theme"
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
	if env == "" {
		env = "database_description"
	}
	tables := ListAllTableNames(conn, dc.Database)
	defer s.FreeResources() // Close connection before freeing resources
	files.Env = env
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
			f.Type = s.TypeMap.TypeConvert[v.DataType]
			f.Name = s.GoFields[i]
			if v.ColumnKey == "PRI" {
				f.PrimaryKey = true
			} else {
				f.PrimaryKey = false
			}
			rl := GetRelationship(v.ColumnName, rt)
			if rl != nil {
				var foreign FieldElements
				if rl.Relationship == MTO && v.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					foreign.Name = StandardizeName(rl.Table)
					foreign.Source = rl.Table
					foreign.Type = "*[]" + StandardizeName(rl.Table)
					foreign.ForeignKey = rl.Column
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
	data, err := json.MarshalIndent(&files, "", " ")
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
	tables := ListAllTableNames(conn, dc.Database)
	defer s.FreeResources() // Close connection before freeing resources
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	s.InitSqlTablesData()
	files.Env = env
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
			f.Type = s.TypeMap.TypeConvert[v.DataType]
			if v.ColumnKey == "PRI" {
				f.PrimaryKey = true
			} else {
				f.PrimaryKey = false
			}
			rl := GetRelationship(v.ColumnName, rt)
			if rl != nil {
				var foreign FieldElements
				if rl.Relationship == MTO && v.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					foreign.Name = StandardizeName(rl.Table)
					foreign.Source = rl.Table
					foreign.Type = "*[]" + StandardizeName(rl.Table)
					foreign.ForeignKey = rl.Column
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
	data, err := json.MarshalIndent(&files, "", " ")
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

func InputUI(dc *DatabaseConfig) fyne.Window {
	app := fApp.New()
	app.Settings().SetTheme(theme.LightTheme())
	w := app.NewWindow("Database Json Generator")
	w.Resize(fyne.Size{
		Width: 640,
	})
	dialectEntry := widget.NewEntry()
	dialectEntry.OnChanged = dc.SetDialect
	usernameEntry := widget.NewEntry()
	usernameEntry.OnChanged = dc.SetUsername
	passwordEntry := widget.NewEntry()
	passwordEntry.OnChanged = dc.SetPassword
	passwordEntry.Password = true
	hostEntry := widget.NewEntry()
	hostEntry.OnChanged = dc.SetHost
	portEntry := widget.NewEntry()
	portEntry.OnChanged = dc.SetPort
	databaseEntry := widget.NewEntry()
	databaseEntry.OnChanged = dc.SetDatabaseName
	w.SetContent(widget.NewVBox(
		widget.NewLabel("Dialect:"),
		dialectEntry,
		widget.NewLabel("User:"),
		usernameEntry,
		widget.NewLabel("Password:"),
		passwordEntry,
		widget.NewLabel("Host:"),
		hostEntry,
		widget.NewLabel("Port:"),
		portEntry,
		widget.NewLabel("Database:"),
		databaseEntry,
		widget.NewButton("Generate Database Json Description", func() {
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
		}),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))
	return w
}
