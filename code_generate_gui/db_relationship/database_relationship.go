package db_relationship

import (
	"context"
	"encoding/json"
	"github.com/go-generator/project"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	. "../cache/yaml/cache_cipher"
	. "../db_config"
	. "./constants"
	. "./database_models"
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/go-generator/database"
	newConfig "github.com/go-generator/database/db_config"
	. "github.com/go-generator/metadata"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sqweek/dialog"
	"golang/code_generate_gui/constants"
	"golang/code_generate_gui/db_relationship/common"
	"golang/code_generate_gui/map_type"
	"golang/code_generate_gui/utils"
	"golang/code_generate_gui/working_directory"
)

var (
	env             = "model"
	JavaTemplateDir = filepath.Join(working_directory.GetWorkingDirectory(), "db_relationship", "template")
)

// 1-1 -> both fields are unique
// 1-n -> only one field is unique
// n-n -> both fields are not unique
// self reference will be in the same table with the same datatype

func DatabaseRelationships(dbConfig *DatabaseConfig, conn *gorm.DB) ([]RelationshipTables, []string) {
	rt := NewRelationshipTables(dbConfig, conn)
	jt := ListAllJoinTablesWithCompositeKey(dbConfig.Database, conn, rt)
	for i := range rt {
		rt[i].Relationship = FindRelationShip(dbConfig.Database, conn, jt, &rt[i])
	}
	joinTable := ListAllJoinTablesWithCompositeKey(dbConfig.Database, conn, rt)
	return rt, joinTable
}

func JsonDescriptionGenerator(env, output string, conn *gorm.DB, dc *DatabaseConfig, rt []RelationshipTables) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var (
		err    error
		files  FilesDetails
		entity []string
	)
	files.Env = []string{"search_model", "config", "controller", "service/impl", "route", "service"}
	var folderOutput Folders
	if env == "" {
		env = "model"
	}
	tables := ListAllTableNames(conn, dc.Database, dc.Dialect)
	files.Model = env
	path := "./" + output + "/"
	fileDirectory := path + output + ".json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.Mkdir(path, 0777)
		if err != nil {
			return err
		}
	}
	for _, v := range tables {
		var (
			sqlTable SqlTablesData
			m        Model
		)
		sqlTable.TypeConvert = map_type.RetrieveTypeMap()
		m.Name = v
		sqlTable.InitSqlTable(dc.Database, v, conn)
		sqlTable.StandardizeFieldsName()
		for i, k := range sqlTable.SqlTable {
			var f Field
			if sqlTable.ContainCompositeKey {
				f.Source = ToLower(sqlTable.GoFields[i])
			} else {
				if k.ColumnKey == "PRI" {
					f.Source = "_id"
				} else {
					f.Source = ToLower(sqlTable.GoFields[i])
				}
			}
			f.Type = sqlTable.TypeConvert[k.DataType]
			f.Name = sqlTable.GoFields[i]
			if k.ColumnKey == "PRI" {
				f.Id = true
			} else {
				f.Id = false
			}
			rl := GetRelationship(k.ColumnName, rt)
			if rl != nil {
				var foreign Field
				if rl.Relationship == ManyToOne && k.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					var relationship Relationship
					relationship.Ref = k.TableName
					foreign.Name = StandardizeName(rl.Table)
					foreign.Source = rl.Table
					foreign.Type = "*[]" + StandardizeName(rl.Table)
					relationship.Fields = append(relationship.Fields, Link{
						Column: rl.Column,
						To:     rl.ReferencedColumn,
					})
					m.Arrays = append(m.Arrays, relationship)
					m.Fields = append(m.Fields, foreign)
				}
			}
			m.Fields = append(m.Fields, f)
		}
		files.Files = append(files.Files, m)
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

func JsonUI(env, filePath string, conn *gorm.DB, dc *DatabaseConfig, rt []RelationshipTables, jt []string, opt bool) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var (
		err    error
		files  FilesDetails
		output Folders
		entity []string
	)
	typeMap := map_type.RetrieveTypeMap()
	files.Env = []string{"search_model", "config", "controller", "service/impl", "route", "service"}
	tables := ListAllTableNames(conn, dc.Database, dc.Dialect)
	files.Model = env
	for _, v := range tables {
		var (
			sqlTable SqlTablesData
			m        Model
		)
		sqlTable.TypeConvert = utils.CopyMap(typeMap)
		if opt && utils.IsContainedInStrings(v, jt) { // Not generate model for many to many tables
			continue
		}
		m.Name = v
		sqlTable.InitSqlTable(dc.Database, v, conn)
		sqlTable.StandardizeFieldsName()
		for i, v := range sqlTable.SqlTable {
			var f Field
			if sqlTable.ContainCompositeKey {
				f.Source = ToLower(sqlTable.GoFields[i])
			} else {
				if v.ColumnKey == "PRI" {
					f.Source = "_id"
				} else {
					f.Source = ToLower(sqlTable.GoFields[i])
				}
			}
			f.Name = sqlTable.GoFields[i]
			f.Type = sqlTable.TypeConvert[v.DataType]
			if v.ColumnKey == "PRI" {
				f.Id = true
			}
			rl := GetRelationship(v.ColumnName, rt)
			if rl != nil {
				var relationship Relationship
				var foreign Field
				foreign.Name = StandardizeName(rl.Table)
				foreign.Source = rl.Table
				foreign.Type = "*[]" + StandardizeName(rl.Table)
				if rl.Relationship == ManyToOne && v.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					relationship.Ref = rl.Table
					relationship.Fields = append(relationship.Fields, Link{
						Column: rl.Column,
						To:     rl.ReferencedColumn,
					})
					if m.Arrays == nil {
						m.Arrays = append(m.Arrays, relationship)
					} else {
						for j := range m.Arrays {
							if m.Arrays[j].Ref == relationship.Ref {
								m.Arrays[j].Fields = append(m.Arrays[j].Fields, relationship.Fields...)
								break
							}
							if j == len(m.Arrays)-1 {
								m.Arrays = append(m.Arrays, relationship)
							}
						}
					}
					for i := range m.Fields {
						if m.Fields[i] == foreign {
							break
						}
						if i == len(m.Fields)-1 {
							m.Fields = append(m.Fields, foreign)
						}
					}
				}
			}
			m.Fields = append(m.Fields, f)
		}
		files.Files = append(files.Files, m)
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

func OldWriteMetadata(app fyne.App, dc *DatabaseConfig, conn *gorm.DB, optimize bool) {
	err := dc.ValidateDatabaseConfig()
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	filename, errFile := dialog.File().Filter("json files", "json").Title("Save As").Save()
	if errFile != nil {
		ShowWindows(app, "Error", errFile.Error())
		return
	}
	rl, jt := DatabaseRelationships(dc, conn)
	err = JsonUI(env, filename+".json", conn, dc, rl, jt, optimize)
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	ShowWindows(app, "Success", "Generated Database Json File Successfully")
}
func SaveMetadataJson(projectStruct interface{}, filePath string) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	data, err := json.MarshalIndent(&projectStruct, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filePath, data, 0644) // Create and write files
	if err != nil {
		return err
	}
	return err
}
func oldSaveMetadataJson(env, filePath string, metaDataList []Model) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var (
		err    error
		files  FilesDetails
		output Folders
		entity []string
	)
	files.Env = []string{"search_model", "config", "controller", "service/impl", "route", "service"}
	files.Model = env
	files.Files = metaDataList
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
func WriteMetadata(app fyne.App, dc *DatabaseConfig, conn *gorm.DB, optimize bool, seperateFile bool) {
	err := dc.ValidateDatabaseConfig()
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	inFile, err := dialog.File().Filter("json file", "json").Title("Open As").Load()
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	outFile, err := dialog.File().Filter("json files", "json").Title("Save As").Save()
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	modelFile := ""
	if seperateFile {
		modelFile, err = dialog.File().Filter("json files", "json").Title("Save Model File As").Save()
		if err != nil {
			ShowWindows(app, "Error", err.Error())
			return
		}
	}

	var t project.DatabaseAdapter
	t = &database.DefaultMetadataService{Config: newConfig.DatabaseConfig{
		Dialect:  dc.Dialect,
		Host:     dc.Host,
		Port:     dc.Port,
		Database: dc.Database,
		User:     dc.User,
		Password: dc.Password,
	}}
	//metadataList := t.ToMetadata(context.Background(), conn)
	var u project.ProjectService
	u = &project.GoMongoProjectService{}
	projectStruct, err := u.CreateProjectByAdapter(context.Background(), inFile, "hotelManagement", t, conn)
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	if seperateFile {
		projectStruct.ModelsFile = modelFile
		err = SaveMetadataJson(projectStruct.Models, modelFile)
		if err != nil {
			ShowWindows(app, "Error", err.Error())
			return
		}
		projectStruct.Models = nil
		err = SaveMetadataJson(projectStruct, outFile)
		if err != nil {
			ShowWindows(app, "Error", err.Error())
			return
		}
	} else {
		err = SaveMetadataJson(projectStruct, outFile)
		if err != nil {
			ShowWindows(app, "Error", err.Error())
			return
		}
	}
	ShowWindows(app, "Success", "Generated Database Json File Successfully")
}

//func OlderWriteMetadata(app fyne.App, dc *DatabaseConfig, conn *gorm.DB, optimize bool) {
//	err := dc.ValidateDatabaseConfig()
//	if err != nil {
//		ShowWindows(app, "Error", err.Error())
//		return
//	}
//
//	filename, err := dialog.File().Filter("json files", "json").Title("Save As").Save()
//	if err != nil {
//		ShowWindows(app, "Error", err.Error())
//		return
//	}
//	var t project.DatabaseAdapter
//	t = &database.DefaultMetadataService{Config: newConfig.DatabaseConfig{
//		Dialect:  dc.Dialect,
//		Host:     dc.Host,
//		Port:     dc.Port,
//		Database: dc.Database,
//		User:     dc.User,
//		Password: dc.Password,
//	}}
//	metadataList := t.ToMetadata(context.Background(), conn)
//	err = SaveMetadataJson(env, filename, metadataList)
//	if err != nil {
//		ShowWindows(app, "Error", err.Error())
//		return
//	}
//	ShowWindows(app, "Success", "Generated Database Json File Successfully")
//}

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
	var opt bool
	var sep bool
	optimizeEntry := widget.NewCheck("Optimization", func(b bool) {
		opt = b
	})
	seperateEntry := widget.NewCheck("Seperate", func(b bool) {
		sep = b
	})
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
	portEntry.Text = strconv.Itoa(dc.Port)
	portEntry.OnChanged = func(s string) {
		if s == "" {
			return
		}
		tmp, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			ShowWindows(app, "Error", constants.ErrInvPort)
			portEntry.SetText(strconv.Itoa(dc.Port))
			return
		}
		if tmp < 1 {
			ShowWindows(app, "Error", constants.ErrInvPort)
			portEntry.SetText(strconv.Itoa(dc.Port))
			return
		}
		dc.Port = int(tmp)
		portEntry.SetText(strconv.Itoa(dc.Port))
	}
	databaseEntry := widget.NewEntry()
	databaseEntry.OnChanged = dc.SetDatabaseName
	databaseEntry.Text = dc.Database
	executeButton := widget.NewButton("Generate Database Json Description", func() {
		conn, err := dc.ConnectToSqlServer()
		if err != nil {
			ShowWindows(app, "Error", err.Error())
			return
		}
		defer func() {
			err = conn.Close()
			if err != nil {
				ShowWindows(app, "Error", err.Error())
			}
		}()
		WriteMetadata(app, dc, conn, opt, sep)
		if temp != *dc {
			err := WriteCacheFile(cache, dc, encryptField)
			if err != nil {
				log.Println(err)
			}
		}
		window.Close()
	})
	providerEntry := widget.NewRadio([]string{"mysql", "postgres", "mssql", "sqlite3"}, func(s string) {
		dc.SetDialect(s)
	})
	providerEntry.Selected = dc.Dialect
	providerEntry.Refresh()
	window.SetContent(widget.NewVBox(
		optimizeEntry,
		seperateEntry,
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
	))
	window.CenterOnScreen()
	return window
}

func JavaUI() error {
	var output Folders
	file, err := dialog.File().Title("Select json").Filter("json file", "json").Load()
	if err != nil {
		return err
	}
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &output)
	if err != nil {
		return err
	}
	outFolder, err := dialog.Directory().Title("Output java files").Browse()
	if err != nil {
		return err
	}
	err = WriteJavaFiles(output, outFolder)
	return err
}

func GetConnection(output *Folders) []Connection {
	var connection []Connection
	for _, v := range output.ModelFile[0].Files {
		if v.Arrays != nil {
			for _, v1 := range v.Arrays {
				tmp := Connection{
					TableName:       v.Name,
					ReferencedTable: v1.Ref,
					Fields:          nil,
				}
				for _, v2 := range v1.Fields {
					field := Link{
						Column: v2.Column,
						To:     v2.To,
					}
					tmp.Fields = append(tmp.Fields, field)
				}
				connection = append(connection, tmp)
			}
		}
	}
	return connection
}

func ParseJavaTemplate(jTmplFilePath, jOutput, tableName string, javaTemplate interface{}) error {
	funcMap := template.FuncMap{
		"ToTitle": strings.Title,
		"ToLower": strings.ToLower,
	}
	tmplName := filepath.Base(jTmplFilePath)
	jTmpl := template.New(tmplName).Funcs(funcMap)
	jT, err := jTmpl.ParseFiles(jTmplFilePath)
	if err != nil {
		return err
	}
	switch v := javaTemplate.(type) {
	case *JavaTemplate:
		outputDir := filepath.Join(jOutput, strings.Title(tableName)+".java")
		jFile, err := os.OpenFile(outputDir, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err = jFile.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		err = jT.Execute(jFile, v)
	case *JavaComPK:
		outputDir := filepath.Join(jOutput, strings.Title(tableName)+"PK.java")
		jFile, err := os.OpenFile(outputDir, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer func() {
			err = jFile.Close()
			if err != nil {
				log.Println(err)
			}
		}()
		err = jT.Execute(jFile, v)
	}
	return err
}

func GetValueColumns(output Folders, table string) []string {
	var res []string
	for _, v := range output.ModelFile[0].Files {
		if StandardizeName(v.Name) == table {
			for _, v1 := range v.Fields {
				if !v1.Id && !common.IsExisted(StandardizeName(v1.Name), output.ModelFile[0].Entity) {
					res = append(res, ToLower(v1.Name))
				}
			}
			break
		}
	}
	return res
}

func GetAllPrimaryKeys(output Folders, table string) []string {
	var res []string
	for _, v := range output.ModelFile[0].Files {
		if StandardizeName(v.Name) == table {
			for _, v1 := range v.Fields {
				if v1.Id {
					res = append(res, ToLower(v1.Name))
				}
			}
			break
		}
	}
	return res
}

func WriteJavaFiles(output Folders, filePath string) error {
	var err error
	jTmplFile := filepath.Join(JavaTemplateDir, "java_template.tmpl")
	jTmplOtmFile := filepath.Join(JavaTemplateDir, "java_template_otm.tmpl")
	tmplPKFile := filepath.Join(JavaTemplateDir, "java_pk.tmpl")
	connection := GetConnection(&output)
	for i := range output.ModelFile[0].Entity {
		fields := GetValueColumns(output, output.ModelFile[0].Entity[i])
		priCols := GetAllPrimaryKeys(output, output.ModelFile[0].Entity[i])
		javaData := NewJavaTemplate(output.ModelFile[0].Model,
			output.ModelFile[0].Entity[i],
			priCols,
			fields,
		)
		for index := range connection {
			if StandardizeName(connection[index].TableName) == StandardizeName(javaData.TableName) {
				tableRef := TableRef{
					Name:        "",
					JoinColumns: nil,
				}
				for index2 := range connection[index].Fields {
					joinColumn := ColumnRef{
						Col:           connection[index].Fields[index2].Column,
						ReferencedCol: connection[index].Fields[index2].To,
					}
					tableRef.JoinColumns = append(tableRef.JoinColumns, joinColumn)
				}
				tableRef.Name = StandardizeName(connection[index].ReferencedTable)
				javaData.TableRef = append(javaData.TableRef, tableRef)
				break
			}
		}
		if len(javaData.IDFields) > 1 {
			data := NewJavaComPK(javaData.Package, javaData.TableName+"PK", javaData.IDFields)
			err = ParseJavaTemplate(tmplPKFile, filePath, javaData.TableName, data)
			if err != nil {
				return err
			}
		}
		if javaData.TableRef != nil {
			err = ParseJavaTemplate(jTmplOtmFile, filePath, javaData.TableName, javaData)
			if err != nil {
				return err
			}
		} else {
			err = ParseJavaTemplate(jTmplFile, filePath, javaData.TableName, javaData)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func GetCompositeKeys(modelJSON Model) []string {
	count := 0
	pri := make([]string, 0)
	for _, v := range modelJSON.Fields {
		if v.Id {
			count++
			pri = append(pri, v.Name)
		}
	}
	if count > 1 {
		return pri
	}
	return nil
}
