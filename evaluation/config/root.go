package config

import (
	ldap "../../ldap-authentication"
	"github.com/common-go/auth"
	app "github.com/common-go/echo"
	"github.com/common-go/mongo"
)

type Root struct {
	Server   app.ServerConfig   `mapstructure:"server"`
	Ldap     ldap.LDAPConfig    `mapstructure:"ldap"`
	Token    auth.TokenConfig   `mapstructure:"token"`
	Mongo    mongo.MongoConfig  `mapstructure:"mongo"`
}
