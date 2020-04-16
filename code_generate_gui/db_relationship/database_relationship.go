package db_relationship

import (
	"encoding/json"
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
			m        ModelJSON
		)
		sqlTable.TypeConvert = map_type.RetrieveTypeMap()
		m.Name = v
		sqlTable.InitSqlTable(dc.Database, v, conn)
		sqlTable.StandardizeFieldsName()
		for i, k := range sqlTable.SqlTable {
			var f FieldElements
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
				f.PrimaryKey = true
			} else {
				f.PrimaryKey = false
			}
			rl := GetRelationship(k.ColumnName, rt)
			if rl != nil {
				var foreign FieldElements
				if rl.Relationship == ManyToOne && k.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					var relationship Relationship
					relationship.Table = k.TableName
					foreign.Name = StandardizeName(rl.Table)
					foreign.Source = rl.Table
					foreign.Type = "*[]" + StandardizeName(rl.Table)
					relationship.Fields = append(relationship.Fields, Field{
						ColumnName:       rl.Column,
						ReferencedColumn: rl.ReferencedColumn,
					})
					m.Relationships = append(m.Relationships, relationship)
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
			m        ModelJSON
		)
		sqlTable.TypeConvert = utils.CopyMap(typeMap)
		if opt && utils.IsContainedInStrings(v, jt) { // Not generate model for many to many tables
			continue
		}
		m.Name = v
		sqlTable.InitSqlTable(dc.Database, v, conn)
		sqlTable.StandardizeFieldsName()
		for i, v := range sqlTable.SqlTable {
			var f FieldElements
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
				f.PrimaryKey = true
			}
			rl := GetRelationship(v.ColumnName, rt)
			if rl != nil {
				var relationship Relationship
				var foreign FieldElements
				foreign.Name = StandardizeName(rl.Table)
				foreign.Source = rl.Table
				foreign.Type = "*[]" + StandardizeName(rl.Table)
				if rl.Relationship == ManyToOne && v.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
					relationship.Table = rl.Table
					relationship.Fields = append(relationship.Fields, Field{
						ColumnName:       rl.Column,
						ReferencedColumn: rl.ReferencedColumn,
					})
					if m.Relationships == nil {
						m.Relationships = append(m.Relationships, relationship)
					} else {
						for j := range m.Relationships {
							if m.Relationships[j].Table == relationship.Table {
								m.Relationships[j].Fields = append(m.Relationships[j].Fields, relationship.Fields...)
								break
							}
							if j == len(m.Relationships)-1 {
								m.Relationships = append(m.Relationships, relationship)
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

func WriteMetadata(app fyne.App, dc *DatabaseConfig, conn *gorm.DB, optimize bool) {
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

//func WriteJavaMetadata(app fyne.App, dc *DatabaseConfig, conn *gorm.DB, optimize bool) {
//	err := dc.ValidateDatabaseConfig()
//	if err != nil {
//		ShowWindows(app, "Error", err.Error())
//		return
//	}
//	javaPath, err := dialog.Directory().Title("Java files path").Browse()
//	if err != nil {
//		ShowWindows(app, "Error", err.Error())
//		return
//	}
//	rl, jt := DatabaseRelationships(dc, conn)
//	err = JavaUI(env, javaPath, conn, dc, rl, jt, optimize)
//	if err != nil {
//		ShowWindows(app, "Error", err.Error())
//		return
//	}
//	ShowWindows(app, "Success", "Generated Database Java Files Successfully")
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
	optimizeEntry := widget.NewCheck("Optimization", func(b bool) {
		opt = b
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
		temp, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			ShowWindows(app, "Error", constants.ErrInvPort)
			portEntry.SetText(strconv.Itoa(dc.Port))
			return
		}
		if temp < 1 {
			ShowWindows(app, "Error", constants.ErrInvPort)
			portEntry.SetText(strconv.Itoa(dc.Port))
			return
		}
		dc.Port = int(temp)
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
		WriteMetadata(app, dc, conn, opt)
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

//func JavaInputUI(dc *DatabaseConfig, app fyne.App, cache, encryptField string) fyne.Window {
//	var temp DatabaseConfig
//	err := ReadCacheFile(cache, &temp, encryptField)
//	if err != nil {
//		log.Println(err)
//	}
//	window := app.NewWindow("Database Java Files Generator")
//	window.Resize(fyne.Size{
//		Width: 640,
//	})
//	var opt bool
//	optimizeEntry := widget.NewCheck("Optimization", func(b bool) {
//		opt = b
//	})
//	usernameEntry := widget.NewEntry()
//	usernameEntry.OnChanged = dc.SetUsername
//	usernameEntry.Text = dc.User
//	passwordEntry := widget.NewEntry()
//	passwordEntry.OnChanged = dc.SetPassword
//	passwordEntry.Text = dc.Password
//	passwordEntry.Password = true
//	hostEntry := widget.NewEntry()
//	hostEntry.OnChanged = dc.SetHost
//	hostEntry.Text = dc.Host
//	portEntry := widget.NewEntry()
//	portEntry.Text = strconv.Itoa(dc.Port)
//	portEntry.OnChanged = func(s string) {
//		if s == "" {
//			return
//		}
//		temp, err := strconv.ParseInt(s, 10, 32)
//		if err != nil {
//			ShowWindows(app, "Error", constants.ErrInvPort)
//			portEntry.SetText(strconv.Itoa(dc.Port))
//			return
//		}
//		if temp < 1 {
//			ShowWindows(app, "Error", constants.ErrInvPort)
//			portEntry.SetText(strconv.Itoa(dc.Port))
//			return
//		}
//		dc.Port = int(temp)
//		portEntry.SetText(strconv.Itoa(dc.Port))
//	}
//	databaseEntry := widget.NewEntry()
//	databaseEntry.OnChanged = dc.SetDatabaseName
//	databaseEntry.Text = dc.Database
//	executeButton := widget.NewButton("Generate Database Java Description", func() {
//		conn, err := dc.ConnectToSqlServer()
//		if err != nil {
//			ShowWindows(app, "Error", err.Error())
//			return
//		}
//		defer func() {
//			err = conn.Close()
//			if err != nil {
//				ShowWindows(app, "Error", err.Error())
//			}
//		}()
//		WriteJavaMetadata(app, dc, conn, opt)
//		if temp != *dc {
//			err := WriteCacheFile(cache, dc, encryptField)
//			if err != nil {
//				log.Println(err)
//			}
//		}
//		window.Close()
//	})
//	providerEntry := widget.NewRadio([]string{"mysql", "postgres", "mssql", "sqlite3"}, func(s string) {
//		dc.SetDialect(s)
//	})
//	providerEntry.Selected = dc.Dialect
//	providerEntry.Refresh()
//	window.SetContent(widget.NewVBox(
//		optimizeEntry,
//		widget.NewLabel("Provider:"),
//		providerEntry,
//		widget.NewLabel("User:"),
//		usernameEntry,
//		widget.NewCheck("Show Password", func(b bool) {
//			if b {
//				passwordEntry.Password = false
//			} else {
//				passwordEntry.Password = true
//			}
//			passwordEntry.Refresh()
//		}),
//		widget.NewLabel("Password:"),
//		passwordEntry,
//		widget.NewLabel("Host:"),
//		hostEntry,
//		widget.NewLabel("Port:"),
//		portEntry,
//		widget.NewLabel("Database:"),
//		databaseEntry,
//		executeButton,
//	))
//	window.CenterOnScreen()
//	return window
//}

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
	//var (
	//	err    error
	//	files  FilesDetails
	//	output Folders
	//	entity []string
	//)
	//typeMap := map_type.RetrieveTypeMap()
	//files.Env = []string{"search_model", "config", "controller", "service/impl", "route", "service"}
	//tables := ListAllTableNames(conn, dc.Database, dc.Dialect)
	//files.Model = env
	//for _, v := range tables {
	//	var (
	//		sqlTable SqlTablesData
	//		m        ModelJSON
	//	)
	//	sqlTable.TypeConvert = utils.CopyMap(typeMap)
	//	if opt && utils.IsContainedInStrings(v, jt) { // Not generate model for many to many tables
	//		continue
	//	}
	//	m.Name = v
	//	sqlTable.InitSqlTable(dc.Database, v, conn)
	//	sqlTable.StandardizeFieldsName()
	//	for i, v := range sqlTable.SqlTable {
	//		var f FieldElements
	//		if sqlTable.ContainCompositeKey {
	//			f.Source = ToLower(sqlTable.GoFields[i])
	//		} else {
	//			if v.ColumnKey == "PRI" {
	//				f.Source = "_id"
	//			} else {
	//				f.Source = ToLower(sqlTable.GoFields[i])
	//			}
	//		}
	//		f.Name = sqlTable.GoFields[i]
	//		f.Type = sqlTable.TypeConvert[v.DataType]
	//		if v.ColumnKey == "PRI" {
	//			f.PrimaryKey = true
	//		} else {
	//			f.PrimaryKey = false
	//		}
	//		rl := GetRelationship(v.ColumnName, rt)
	//		if rl != nil {
	//			log.Println(rl)
	//			var relationship Relationship
	//			var foreign FieldElements
	//			if rl.Relationship == ManyToOne && v.TableName == rl.ReferencedTable { // Have Many to One relation, add a field to the current struct
	//				relationship.Table = rl.Table
	//				foreign.Name = StandardizeName(rl.Table)
	//				foreign.Source = rl.Table
	//				foreign.Type = "*[]" + StandardizeName(rl.Table)
	//				relationship.Fields = append(relationship.Fields, Field{
	//					ColumnName:       rl.Column,
	//					ReferencedColumn: rl.ReferencedColumn,
	//				})
	//				m.Relationships = append(m.Relationships, relationship)
	//				m.Fields = append(m.Fields, foreign)
	//			}
	//		}
	//		m.Fields = append(m.Fields, f)
	//	}
	//	files.Files = append(files.Files, m)
	//}
	//for _, v := range files.Files {
	//	entity = append(entity, StandardizeName(v.Name))
	//}
	//files.Entity = entity
	//output.ModelFile = append(output.ModelFile, files)
	err = WriteJavaFiles(output, outFolder)
	//TODO: Write Java files from output
	//data, err := json.MarshalIndent(&output, "", " ")
	//if err != nil {
	//	return err
	//}
	//err = ioutil.WriteFile(filePath, data, 0644) // Create and write files
	return err
}

//TODO: Write Java files

func WriteJavaFiles(output Folders, filePath string) error {
	var connection []Connection
	//tables := ListAllTableNames(conn, dc.Database, dc.Dialect)
	//rl, _ := DatabaseRelationships(dc, conn)
	for _, v := range output.ModelFile[0].Files {
		if v.Relationships != nil {
			for _, v1 := range v.Relationships {
				tmp := Connection{
					TableName:       v.Name,
					ReferencedTable: v1.Table,
					Fields:          nil,
				}
				for _, v2 := range v1.Fields {
					field := Field{
						ColumnName:       v2.ColumnName,
						ReferencedColumn: v2.ReferencedColumn,
					}
					tmp.Fields = append(tmp.Fields, field)
				}
				connection = append(connection, tmp)
			}
		}
	}
	for i := range output.ModelFile[0].Entity {
		fields := GetValueColumns(output, output.ModelFile[0].Entity[i])
		priCols := GetAllPrimaryKeys(output, output.ModelFile[0].Entity[i])
		javaData := NewJavaTemplate(output.ModelFile[0].Model, output.ModelFile[0].Entity[i],
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
						Col:           connection[index].Fields[index2].ColumnName,
						ReferencedCol: connection[index].Fields[index2].ReferencedColumn,
					}
					tableRef.JoinColumns = append(tableRef.JoinColumns, joinColumn)
				}
				tableRef.Name = StandardizeName(connection[index].ReferencedTable)
				javaData.TableRef = append(javaData.TableRef, tableRef)
				break
			}
		}
		jTmplFile := filepath.Join(JavaTemplateDir, "java_template.tmpl")
		jTmplOtmFile := filepath.Join(JavaTemplateDir, "java_template_otm.tmpl")
		tmplPKFile := filepath.Join(JavaTemplateDir, "java_pk.tmpl")
		if len(javaData.IDFields) > 1 {
			data := NewJavaComPK(javaData.Package, javaData.TableName+"PK", javaData.IDFields)
			err := ParseJavaTemplate(tmplPKFile, filePath, javaData.TableName, data)
			if err != nil {
				return err
			}
		}
		if javaData.TableRef != nil {
			err := ParseJavaTemplate(jTmplOtmFile, filePath, javaData.TableName, javaData)
			if err != nil {
				return err
			}
		} else {
			err := ParseJavaTemplate(jTmplFile, filePath, javaData.TableName, javaData)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func GetCompositeKeys(modelJSON ModelJSON) []string {
	count := 0
	pri := make([]string, 0)
	for _, v := range modelJSON.Fields {
		if v.PrimaryKey {
			count++
			pri = append(pri, v.Name)
		}
	}
	if count > 1 {
		return pri
	}
	return nil
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
				if !v1.PrimaryKey && !common.IsExisted(StandardizeName(v1.Name), output.ModelFile[0].Entity) {
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
				if v1.PrimaryKey {
					res = append(res, ToLower(v1.Name))
				}
			}
			break
		}
	}
	return res
}
