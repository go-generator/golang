package database_models

import (
	. "../../db_config"
	. "../common"
	. "../constants"
	"github.com/jinzhu/gorm"
)

type TestMap struct {
	All map[string]interface{}
}

type RelationshipTables struct {
	Table            string `gorm:"column:TABLE_NAME"`
	Column           string `gorm:"column:COLUMN_NAME"`
	ReferencedTable  string `gorm:"column:REFERENCED_TABLE_NAME"`
	ReferencedColumn string `gorm:"column:REFERENCED_COLUMN_NAME"`
	Relationship     string
}

type LowercaseRelationshipTables struct {
	Table            string `gorm:"column:table_name"`
	Column           string `gorm:"column:column_name"`
	ReferencedTable  string `gorm:"column:referenced_table_name"`
	ReferencedColumn string `gorm:"column:referenced_column_name"`
	Relationship     string
}

func checkRelation(check, checkReference bool, database string, connection *gorm.DB, rt *RelationshipTables) string {
	// Already cover the ManyToMany case where a joined table consists of two or more primary key tags that are all foreign keys
	isPrimaryTag := CheckPrimaryTag(database, rt.Table, rt.Column, connection)
	isReferencedPrimaryTag := CheckPrimaryTag(database, rt.ReferencedTable, rt.ReferencedColumn, connection)
	if !checkReference {
		return UNSUPPORTED
	}
	if check {
		count := CompositeKeyColumns(connection, database, rt.Table)
		if len(count) == 1 { // Only one column has Primary Tag
			if isPrimaryTag && isReferencedPrimaryTag { // Both are Primary key
				return OneToOne
			}
			if !isPrimaryTag && isReferencedPrimaryTag { // Column is only a foreign key referenced to other primary key
				return ManyToOne
			}
		}
		if len(count) > 1 { // Consist of at least one column that has primary key tag and not referenced to other table
			return OneToMany
		}
	}
	if !check {
		return ManyToOne
	}
	if !check && !checkReference {
		return UNSUPPORTED
	}
	return UNKNOWN
}

func FindRelationShip(database string, connection *gorm.DB, joinedTable []string, rt *RelationshipTables) string {
	check := CheckUniqueness(database, rt.Table, rt.Column, connection)
	checkReference := CheckUniqueness(database, rt.ReferencedTable, rt.ReferencedColumn, connection)
	for _, v := range joinedTable {
		if rt.Table == v {
			return ManyToMany
		}
	}
	return checkRelation(check, checkReference, database, connection, rt)
}

func CheckForeignKey(table, column string, rt []RelationshipTables) bool {
	for _, v := range rt {
		if v.Table == table && v.Column == column {
			return true
		}
	}
	return false
} // Check if the column of the table is a foreign key

func IsJoinedTable(table string, columns []string, rt []RelationshipTables) bool {
	for _, v := range columns {
		if CheckForeignKey(table, v, rt) == false {
			return false
		}
	}
	return true
}

func NewRelationshipTables(dbConfig *DatabaseConfig, connection *gorm.DB) []RelationshipTables { // mysql only for now
	var res []RelationshipTables
	var lres []LowercaseRelationshipTables
	switch dbConfig.Dialect {
	case "postgres":
		connection.Raw("SELECT tc.table_schema, tc.constraint_name, tc.table_name AS TABLE_NAME, kcu.column_name AS COLUMN_NAME, ccu.table_schema AS foreign_table_schema,  ccu.table_name AS REFERENCED_TABLE_NAME, ccu.column_name AS REFERENCED_COLUMN_NAME FROM information_schema.table_constraints AS tc JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name AND tc.table_schema = kcu.table_schema JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name AND ccu.table_schema = tc.table_schema WHERE tc.constraint_type = 'FOREIGN KEY';").Scan(&lres)
	case "mysql":
		connection.Table("information_schema.key_column_usage").Select("*").Where("constraint_schema = '" + dbConfig.Database + "' and referenced_table_schema is not null and referenced_table_name is not null and referenced_column_name is not null").Scan(&res)
	}
	return res
} // Find all columns, table and its referenced columns, tables

func ListAllJoinTablesWithCompositeKey(database string, conn *gorm.DB, rt []RelationshipTables) []string {
	var joinTable []string
	tables := ListAllTableNames(conn, database, conn.Dialect().GetName())
	for _, v := range tables {
		columns := ContainCompositeKey(database, v, conn)
		if len(columns) > 1 && IsJoinedTable(v, GetCompositeColumnName(columns), rt) {
			joinTable = append(joinTable, v)
		}
	}
	return joinTable
}

func GetRelationship(column string, rt []RelationshipTables) *RelationshipTables {
	for _, v := range rt {
		if column == v.ReferencedColumn {
			return &v
		}
	}
	return nil
}
