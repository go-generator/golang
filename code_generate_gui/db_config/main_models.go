package db_config

import (
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type FilesDetails struct {
	Env    []string    `json:"env"`
	Entity []string    `json:"entity"`
	Model  string      `json:"model"`
	Files  []ModelJSON `json:"files"`
}

type Folders struct {
	ModelFile []FilesDetails `json:"folders"`
}

type ModelJSON struct {
	Name       string          `json:"name"`
	Source     string          `json:"source"`
	ConstValue []Const         `json:"const"`
	TypeAlias  []TypeAlias     `json:"type_alias"`
	Fields     []FieldElements `json:"fields"`
	WriteFile  strings.Builder `json:"-"`
}

type Const struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type TypeAlias struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type FieldElements struct {
	Name          string         `json:"name"`
	Source        string         `json:"source"`
	Type          string         `json:"type"`
	ForeignKey    string         `json:"foreignKey"`
	Relationships []Relationship `json:"relationships"`
	PrimaryKey    bool           `json:"primaryKey"`
}

type Relationship struct {
	ReType string
	Ref    References
}

type References struct {
	Table   string
	RefCols []string
}
