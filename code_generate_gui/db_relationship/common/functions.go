package common

import (
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
)

type MySQLUniqueness struct {
	ColumnName string `gorm:"column:Column_name"`
	NonUnique  bool   `gorm:"column:Non_unique"` // False mean it's unique, True means it can contain duplicate
	KeyName    string `gorm:"column:Key_name"`
}

type PostgresUnique struct {
	// TODO: Add struct
	TableName string `gorm:"column:tablename"`
	IndexName string `gorm:"column:indexname"`
	IndexDef  string `gorm:"column:indexdef"`
}

type ColumnName struct {
	ColumnName string `gorm:"column:COLUMN_NAME"`
}

func CheckUniqueness(database, table, column string, conn *gorm.DB) bool {
	var (
		mySqlIndex    []MySQLUniqueness
		postgresIndex []PostgresUnique
	)
	//log.Println(conn.Dialect().GetName())
	switch conn.Dialect().GetName() {
	case "postgres":
		//TODO: Check uniqueness for Postgres database
		postgresString := "SELECT * FROM pg_indexes WHERE tablename = '" + table + "'"
		conn.Raw(postgresString).Scan(&postgresIndex)
		for _, v := range postgresIndex {
			if strings.Contains(v.IndexName, "unq") {
				tokens := strings.Split(v.IndexName, "_")
				for i := range tokens {
					if tokens[i] == "unq" {
						columnName := strings.Join(tokens[i:], "_")
						if column == columnName {
							return true
						}
					}
				}
			}
		}
	case "mysql":
		mysqlString := "show indexes from " + database + "." + table
		conn.Raw(mysqlString).Scan(&mySqlIndex)
		for _, v := range mySqlIndex {
			if v.ColumnName == column {
				if v.NonUnique == false {
					return true
				}
			}
		}
	}
	return false
} // Check if a column is unique

func CheckPrimaryTag(database, table, column string, conn *gorm.DB) bool {
	var index []MySQLUniqueness
	sqlString := "show indexes from " + database + "." + table
	conn.Raw(sqlString).Scan(&index)
	for _, v := range index {
		if v.ColumnName == column {
			if v.KeyName == "PRIMARY" {
				return true
			}
		}
	}
	return false
} // Check if a column has primary tag

func ContainCompositeKey(database, table string, conn *gorm.DB) []ColumnName { // Return a slice of ColumnName of the composite key
	var res []ColumnName
	sqlString := "select * from information_schema.KEY_COLUMN_USAGE where table_schema='" + database + "' and table_name ='" + table + "' and constraint_name = 'PRIMARY';"
	conn.Raw(sqlString).Scan(&res)
	return res
}

func GetCompositeColumnName(cn []ColumnName) []string { // Get Column Name of the composite key
	var res []string
	for _, v := range cn {
		res = append(res, v.ColumnName)
	}
	return res
}

func CompositeKeyColumns(conn *gorm.DB, databaseName, table string) []string {
	sqlString := strings.Builder{}
	res := make([]ColumnName, 0)
	sqlString.WriteString("SELECT K.COLUMN_NAME FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS C JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS K ON C.TABLE_NAME = K.TABLE_NAME ")
	sqlString.WriteString("AND C.CONSTRAINT_CATALOG = K.CONSTRAINT_CATALOG ")
	sqlString.WriteString("AND C.CONSTRAINT_SCHEMA = K.CONSTRAINT_SCHEMA ")
	sqlString.WriteString("AND C.CONSTRAINT_NAME = K.CONSTRAINT_NAME ")
	sqlString.WriteString("WHERE C.TABLE_SCHEMA = '" + databaseName + "' AND K.TABLE_NAME='" + table + "' AND  C.CONSTRAINT_TYPE = 'PRIMARY KEY'")
	conn.Raw(sqlString.String()).Scan(&res)
	return GetCompositeColumnName(res)
}

func GetAllStructFields(v interface{}) []string {
	var res []string
	val := reflect.Indirect(reflect.ValueOf(v))
	for i := 0; i < val.NumField(); i++ {
		res = append(res, val.Type().Field(i).Name)
	}
	return res
}

func ShellExecutor(program string, arguments []string) ([]byte, error) {
	cmd := exec.Command(program, arguments...)
	return cmd.Output()
}

func MakeDirectory(path string) (err error) {
	err = os.MkdirAll(path, 0644)
	if err != nil && os.IsNotExist(err) {
		return
	}
	return
}

func IsExisted(v string, s []string) bool {
	for i := range s {
		if v == s[i] {
			return true
		}
	}
	return false
}
