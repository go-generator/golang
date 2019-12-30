package config

import (
	app "github.com/common-go/echo"
	"github.com/common-go/mongo"
)

type Root struct {
	Server app.ServerConfig  `mapstructure:"server"`
	Mongo  mongo.MongoConfig `mapstructure:"mongo"`
}
