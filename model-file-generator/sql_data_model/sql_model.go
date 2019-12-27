package sql_data_model

import (
	"log"
	"strings"

	. "../model_files_functions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

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
	PrimaryKey bool   `json:"primaryKey"`
}

type DatabaseCredentials struct {
	Username     string
	Host         string
	Password     string
	DatabaseName string
}

func (credentials *DatabaseCredentials) ConnectToSqlServer() (*gorm.DB, error) {
	conn, err := gorm.Open("mysql", credentials.Username+":"+credentials.Password+"@("+credentials.Host+")/"+credentials.DatabaseName+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		return nil, err
	}
	return conn, err
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
