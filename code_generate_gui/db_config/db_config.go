package db_config

import (
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"
)

type DatabaseConfig struct {
	Dialect  string `mapstructure:"dialect"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
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
