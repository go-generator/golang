package db_config

import (
	"strings"

	"github.com/jinzhu/gorm"
)

type SqlTable struct {
	TableName              string `gorm:"column:TABLE_NAME"`
	ColumnName             string `gorm:"column:COLUMN_NAME"`
	DataType               string `gorm:"column:DATA_TYPE"`
	IsNullable             string `gorm:"column:IS_NULLABLE"`
	ColumnKey              string `gorm:"column:COLUMN_KEY"`
	CharacterMaximumLength string `gorm:"column:CHARACTER_MAXIMUM_LENGTH"`
}

type SqlTablesData struct {
	TypeConvert         map[string]string `yaml:"typeConvert"`
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

func (s *SqlTablesData) ResetData() {
	s.SqlTable = nil
	s.StructName = ""
	s.GoFields = nil
	s.WriteFile.Reset()
}

func (s *SqlTablesData) InitSqlTable(database string, tableName string, conn *gorm.DB) {
	s.StructName = strings.Title(StandardizeName(tableName))
	conn.Table("information_schema.columns").Select("*").Where("TABLE_SCHEMA= '" + database + "' AND TABLE_NAME = '" + tableName + "'").Scan(&s.SqlTable)
}

func (s *SqlTablesData) WritePackage(packageName string) {
	s.WriteFile.WriteString("package " + packageName + "\n\n")
}

func (s *SqlTablesData) WriteStruct() {
	s.WriteFile.WriteString("type " + s.StructName + " struct {\n")
	//if !s.ContainCompositeKey { // Only one Primary key
	//	for i, v := range s.GoFields {
	//		if s.SqlTable[i].ColumnKey == "PRI" {
	//			s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeConvert.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag("_id") + AddGormTag(s.SqlTable[i].ColumnName, true))
	//			continue
	//		}
	//		s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeConvert.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGormTag(s.SqlTable[i].ColumnName, false))
	//	}
	//} else { // Composite key
	for i, v := range s.GoFields {
		if s.SqlTable[i].ColumnKey == "PRI" {
			s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGormTag(s.SqlTable[i].ColumnName, true))
			continue
		}
		s.WriteFile.WriteString("\t" + AddStructFieldName(v) + "\t" + s.TypeConvert[s.SqlTable[i].DataType] + "\t" + AddJSONTag(ToLower(v)) + AddBSONTag(ToLower(v)) + AddGormTag(s.SqlTable[i].ColumnName, false))
	}
	//}
	s.WriteFile.WriteString("}")
}

func (s *SqlTablesData) CreateContent(packageName string) {
	s.WritePackage(packageName)
	s.WriteStruct()
}
