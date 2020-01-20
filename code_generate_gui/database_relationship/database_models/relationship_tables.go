package database_models

import (
	. "../../database_config"
	. "../common"
	. "../constants"
	"github.com/jinzhu/gorm"
)

type RelationshipTables struct {
	Table            string `gorm:"column:TABLE_NAME"`
	Column           string `gorm:"column:COLUMN_NAME"`
	ReferencedTable  string `gorm:"column:REFERENCED_TABLE_NAME"`
	ReferencedColumn string `gorm:"column:REFERENCED_COLUMN_NAME"`
	Relationship     string
}

func checkRelation(check, checkReference bool, database string, connection *DatabaseConnection, rt *RelationshipTables) string {
	// Already cover the MTM case where a joined table consists of two or more primary key tags that are all foreign keys
	isPrimaryTag := CheckPrimaryTag(database, rt.Table, rt.Column, connection.GetConnection())
	isReferencedPrimaryTag := CheckPrimaryTag(database, rt.ReferencedTable, rt.ReferencedColumn, connection.GetConnection())
	if !checkReference {
		return UNS
	}
	if check {
		count := CompositeKeyColumns(connection.GetConnection(), database, rt.Table)
		if len(count) == 1 { // Only one column has Primary Tag
			if isPrimaryTag && isReferencedPrimaryTag { // Both are Primary key
				return OTO
			}
			if !isPrimaryTag && isReferencedPrimaryTag { // Column is only a foreign key referenced to other primary key
				return MTO
			}
		}
		if len(count) > 1 { // Consist of at least one column that has primary key tag and not referenced to other table
			return OTM
		}
	}
	if !check {
		return MTO
	}
	if !check && !checkReference {
		return UNS
	}
	return UNK
}

func FindRelationShip(database string, connection *DatabaseConnection, joinedTable []string, rt *RelationshipTables) string {
	check := CheckUniqueness(database, rt.Table, rt.Column, connection.GetConnection())
	checkReference := CheckUniqueness(database, rt.ReferencedTable, rt.ReferencedColumn, connection.GetConnection())
	for _, v := range joinedTable {
		if rt.Table == v {
			return MTM
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

func NewRelationshipTables(dbConfig *DatabaseConfig, connection *DatabaseConnection) []RelationshipTables {
	var res []RelationshipTables
	connection.GetConnection().Table("information_schema.key_column_usage").Select("*").Where("constraint_schema = '" + dbConfig.Database + "' and referenced_table_schema is not null and referenced_table_name is not null and referenced_column_name is not null").Scan(&res)
	return res
} // Find all columns, table and its referenced columns, tables

func ListAllJoinTablesWithCompositeKey(database string, conn *gorm.DB, rt []RelationshipTables) []string {
	var joinTable []string
	tables := ListAllTableNames(conn, database)
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
