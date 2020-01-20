package database_models

import (
	"strconv"

	. "../../database_config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

type DatabaseConnection struct {
	connection *gorm.DB
}

func (dbConnection *DatabaseConnection) SetConnection(conn *gorm.DB) {
	dbConnection.connection = conn
}

func NewDatabaseConnection(dc *DatabaseConfig) (*DatabaseConnection, error) {
	var dbConnection DatabaseConnection
	var err error
	port := strconv.Itoa(dc.Port)
	switch dc.Dialect {
	case "mysql":
		dbConnection.connection, err = gorm.Open("mysql", dc.User+":"+dc.Password+"@("+dc.Host+":"+port+")/"+dc.Database+"?charset=utf8&parseTime=True&loc=Local")
	case "postgres":
		dbConnection.connection, err = gorm.Open(dc.Dialect, "user="+dc.User+" dbname="+dc.Database+" password="+dc.Password+" host="+dc.Host+" port="+port+" sslmode=disable")
	case "mssql":
		dbConnection.connection, err = gorm.Open(dc.Dialect, "sqlserver://"+dc.User+":"+dc.Password+"@"+dc.Host+":"+port+"?Database="+dc.Database)
	case "sqlite3":
		dbConnection.connection, err = gorm.Open("sqlite3", dc.Host)
	default:
		dbConnection.connection = nil
		err = errors.New("Incorrect Dialect")
	}
	return &dbConnection, err
}

func (dbConnection *DatabaseConnection) GetConnection() *gorm.DB {
	return dbConnection.connection
}

func (dbConnection *DatabaseConnection) CloseConnection() error {
	return dbConnection.connection.Close()
}
