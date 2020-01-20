package database_config

import (
	"log"
	"strconv"
	"strings"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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
	Env    []string    `json:"env"`
	Entity []string    `json:"entity"`
	Model  string      `json:"model"`
	Files  []ModelJSON `json:"files"`
}

type Folders struct {
	ModelFile []FilesDetails `json:"folders"`
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

type TypeMap struct {
	TypeConvert map[string]string `yaml:"typeConvert"`
}

type SqlTable struct {
	TableName              string `gorm:"column:TABLE_NAME"`
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
	viper.SetConfigName("data_type")
	viper.AddConfigPath("../data_type")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error while reading config file, " + err.Error())
	}
	err := viper.Unmarshal(&s.TypeMap)
	log.Println(s.TypeMap.TypeConvert)
	if err != nil {
		log.Println("Error while unmarshal file, " + err.Error())
	}
}

func (s *SqlTablesData) InitSqlTable(database string, tableName string, conn *gorm.DB) {
	s.StructName = strings.Title(StandardizeName(tableName))
	conn.Table("information_schema.columns").Select("*").Where("TABLE_SCHEMA= '" + database + "' AND TABLE_NAME = '" + tableName + "'").Scan(&s.SqlTable)
}

func (s *SqlTablesData) FreeResources() {
	for k, _ := range s.TypeMap.TypeConvert {
		delete(s.TypeMap.TypeConvert, k)
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
	//if !s.ContainCompositeKey { // Only one Primary key
	//	for i, v := range s.GoFields {
	//		if s.SqlTable[i].ColumnKey == "PRI" {
	//			s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag("_id") + AddGORMTag(s.SqlTable[i].ColumnName, true))
	//			continue
	//		}
	//		s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGORMTag(s.SqlTable[i].ColumnName, false))
	//	}
	//} else { // Composite key
	for i, v := range s.GoFields {
		if s.SqlTable[i].ColumnKey == "PRI" {
			s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGORMTag(s.SqlTable[i].ColumnName, true))
			continue
		}
		s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeMap.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGORMTag(s.SqlTable[i].ColumnName, false))
	}
	//}
	s.WriteFile.WriteString("}")
}

func (s *SqlTablesData) CreateContent(packageName string) {
	s.WritePackage(packageName)
	s.WriteStruct()
}
