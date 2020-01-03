package database_relationship

import (
	"log"

	. "../database_config"
	. "./database_models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 1-1 -> both fields are unique
// 1-n -> only one field is unique
// n-n -> both fields are not unique
// self reference will be in the same table with the same datatype

func DatabaseRelationships(dbConfig DatabaseConfig, conn *gorm.DB) ([]RelationshipTables, []string) {
	var dbConnect DatabaseConnection
	dbConnect.SetConnection(conn)
	rt := NewRelationshipTables(&dbConfig, &dbConnect)
	jt := ListAllJoinTablesWithCompositeKey(dbConfig.Database, dbConnect.GetConnection(), rt)
	for _, v := range rt {
		log.Println(v, " -> ", v.FindRelationShip(dbConfig.Database, &dbConnect, jt))
	}
	joinTable := ListAllJoinTablesWithCompositeKey(dbConfig.Database, dbConnect.GetConnection(), rt)
	return rt, joinTable
}
