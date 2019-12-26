package sql_data_model

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	. "../model_files_functions"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

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

func ModelFilesGenerator(s *SqlTablesData, conn *gorm.DB, tables []string, packageName string) {
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
