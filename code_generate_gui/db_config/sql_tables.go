package db_config

import (
	"log"
	"strings"

	. "../map_type"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

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
	viper.AddConfigPath(DTypeAbsPath)
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
	for k := range s.TypeMap.TypeConvert {
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
