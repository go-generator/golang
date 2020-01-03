package database_config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne"
	fApp "fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/sqweek/dialog"
)

type DatabaseConfig struct {
	Dialect  string `mapstructure:"dialect"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type FilesDetails struct {
	Env   string      `json:"env"`
	Files []ModelJSON `json:"files"`
}

type ModelJSON struct {
	Name       string          `json:"name"`
	Source     string          `json:"source"`
	ConstValue []Const         `json:"const"`
	TypeAlias  []TypeAlias     `json:"type_alias"`
	Fields     []FieldElements `json:"fields"`
	WriteFile  strings.Builder `json:"-"`
}

type Const struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TypeAlias struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FieldElements struct {
	Name       string `json:"name"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	ForeignKey string `json:"foreignKey"`
	PrimaryKey bool   `json:"primaryKey"`
}

func (dc *DatabaseConfig) ConnectToSqlServer() (*gorm.DB, error) {
	var conn *gorm.DB
	var err error
	port := strconv.Itoa(dc.Port)
	switch dc.Dialect {
	case "mysql":
		conn, err = gorm.Open("mysql", dc.User+":"+dc.Password+"@("+dc.Host+":"+port+")/"+dc.Database+"?charset=utf8&parseTime=True&loc=Local")
	case "postgres":
		conn, err = gorm.Open(dc.Dialect, "user="+dc.User+" dbname="+dc.Database+" password="+dc.Password+" host="+dc.Host+" port="+port+" sslmode=disable")
	case "mssql":
		conn, err = gorm.Open(dc.Dialect, "sqlserver://"+dc.User+":"+dc.Password+"@"+dc.Host+":"+port+"?Database="+dc.Database)
	case "sqlite3":
		conn, err = gorm.Open("sqlite3", dc.Host)
	default:
		conn = nil
		err = errors.New("Incorrect Dialect")
	}
	return conn, err
}

func (dc *DatabaseConfig) SetDialect(value string) {
	dc.Dialect = value
}

func (dc *DatabaseConfig) SetUsername(value string) {
	dc.User = value
}

func (dc *DatabaseConfig) SetPassword(value string) {
	dc.Password = value
}

func (dc *DatabaseConfig) SetHost(value string) {
	dc.Host = value
}

func (dc *DatabaseConfig) SetPort(value string) {
	var err error
	dc.Port, err = strconv.Atoi(value)
	if err != nil {
		log.Println(err)
	}
}

func (dc *DatabaseConfig) SetDatabaseName(value string) {
	dc.Database = value
}

func (dc *DatabaseConfig) InputValidation(app fyne.App) {
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
	err := dc.JsonUI(tokens[len(tokens)-2], filename+".json")
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return
	}
	ShowWindows(app, "Success", "Generated Database Json Description Successfully")
}

func (dc *DatabaseConfig) InputUI() fyne.Window {
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
			dc.InputValidation(app)
		}),
		widget.NewButton("Quit", func() {
			app.Quit()
		}),
	))
	return w
}

func ShowWindows(app fyne.App, title, message string) {
	wa := app.NewWindow(title)
	wa.Resize(fyne.Size{
		Width: 320,
	})
	wa.SetContent(widget.NewVBox(
		widget.NewLabel(message),
	))
	wa.Show()
}

func (dc *DatabaseConfig) JsonUI(env, filePath string) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
	var s SqlTablesData
	var files FilesDetails
	conn, err := dc.ConnectToSqlServer()
	log.Println("Connecting to sql server...")
	if err != nil {
		return err
	}
	log.Println("Connection to sql server is established successfully")
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
		return err
	}
	err = ioutil.WriteFile(filePath, data, 0644) // Create and write files
	if err != nil {
		return err
	}
	return err
}

type TypeMap map[string]string

type SqlTable struct {
	ColumnName             string `gorm:"column:COLUMN_NAME"`
	DataType               string `gorm:"column:DATA_TYPE"`
	IsNullable             string `gorm:"column:IS_NULLABLE"`
	ColumnKey              string `gorm:"column:COLUMN_KEY"`
	CharacterMaximumLength string `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
}

type SqlTablesData struct {
	TypeMap             TypeMap
	SqlTable            []SqlTable
	GoFields            []string
	StructName          string
	WriteFile           strings.Builder
	ContainCompositeKey bool
}

func (s *SqlTablesData) StandardizeFieldsName() {
	var count int
	for _, v := range s.SqlTable {
		if v.ColumnKey == "PRI" {
			count++
		}
		s.GoFields = append(s.GoFields, StandardizeName(v.ColumnName))
	}
	if count < 2 {
		s.ContainCompositeKey = false
	} else {
		s.ContainCompositeKey = true
	}
}

func (s *SqlTablesData) InitSqlTablesData() {
	s.InitTypeMap()
}

func (s *SqlTablesData) ResetData() {
	s.SqlTable = nil
	s.StructName = ""
	s.GoFields = nil
	s.WriteFile.Reset()
}

func (s *SqlTablesData) InitTypeMap() {
	s.TypeMap = make(TypeMap)
	s.TypeMap["decimal"] = "int"
	s.TypeMap["char"] = "string"
	s.TypeMap["varchar"] = "string"
	s.TypeMap["int"] = "int"
	s.TypeMap["timestamp"] = "time.Time"
	s.TypeMap["date"] = "time.Time"
	s.TypeMap["text"] = "string"
	s.TypeMap["smallint"] = "int8"
	s.TypeMap["bigint"] = "int64"
}

func (s *SqlTablesData) InitSqlTable(tableName string, conn *gorm.DB) {
	s.StructName = strings.Title(StandardizeName(tableName))
	conn.Table("information_schema.columns").Select("*").Where("TABLE_NAME = '" + tableName + "'").Scan(&s.SqlTable)
}

func (s *SqlTablesData) FreeResources() {
	for k, _ := range s.TypeMap {
		delete(s.TypeMap, k)
	}
	log.Println("Resources are freed successfully")
}

func (s *SqlTablesData) WritePackage(packageName string) string {
	s.WriteFile.WriteString("package " + packageName + "\n\n")
	for _, v := range s.SqlTable {
		if v.DataType == "date" || v.DataType == "timestamp" {
			s.WriteFile.WriteString("import \"time\"\n\n")
			break
		}
	}
	return "package " + packageName + "\n\n"
}

func (s *SqlTablesData) WriteStruct() {
	s.WriteFile.WriteString("type " + s.StructName + " struct {\n")
	if !s.ContainCompositeKey { // Only one Primary key
		for i, v := range s.GoFields {
			if s.SqlTable[i].ColumnKey == "PRI" {
				s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag("_id") + AddGORMTag(s.SqlTable[i].ColumnName, true))
				continue
			}
			s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGORMTag(s.SqlTable[i].ColumnName, false))
		}
	} else { // Composite key
		for i, v := range s.GoFields {
			if s.SqlTable[i].ColumnKey == "PRI" {
				s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGORMTag(s.SqlTable[i].ColumnName, true))
				continue
			}
			s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGORMTag(s.SqlTable[i].ColumnName, false))
		}
	}
	s.WriteFile.WriteString("}")
}

func (s *SqlTablesData) CreateContent(packageName string) {
	s.WritePackage(packageName)
	s.WriteStruct()
}

func (dc *DatabaseConfig) JsonDescriptionGenerator(env, output string, conn *gorm.DB) error { //s *SqlTablesData, conn *gorm.DB, tables []string, packageName, output string) {
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
		return err
	}
	err = ioutil.WriteFile(fileDirectory, data, 0644) // Create and write files
	if err != nil {
		return err
	}
	return err
}
