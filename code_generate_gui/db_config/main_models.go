package db_config

import (
	"strings"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type FilesDetails struct {
	Env    []string `json:"env"`
	Entity []string `json:"entity"`
	Model  string   `json:"model"`
	Files  []Model  `json:"files"`
}

type Folders struct {
	ModelFile []FilesDetails `json:"folders"`
}

type Model struct {
	Name       string          `json:"name"`
	Source     string          `json:"source"`
	ConstValue []Const         `json:"const"`
	TypeAlias  []TypeAlias     `json:"type_alias"`
	Models     []Relationship  `json:"models"`
	Arrays     []Relationship  `json:"arrays"`
	Fields     []Field         `json:"fields"`
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

type Field struct {
	Name       string `json:"name,omitempty"`
	Source     string `json:"source,omitempty"`
	Type       string `json:"type,omitempty"`
	PrimaryKey bool   `json:"primaryKey,omitempty"`
}

type Relationship struct {
	Ref    string `json:"table"`
	Fields []Link `json:"fields"`
}

func (m *Model) ExtractFieldType() map[string]string {
	res := make(map[string]string)
	for _, v := range m.Fields {
		res[v.Name] = v.Type
	}
	return res
}

type Connection struct {
	TableName       string
	ReferencedTable string
	Fields          []Link
}

type Link struct {
	Column string `json:"column"`
	Ref    string `json:"ref"`
}

type JavaComPK struct {
	Package    string
	PKName     string
	Array      []string
	Capitalize []string
	AllGet     string
}

func NewJavaComPK(pkg, pkName string, arr []string) *JavaComPK {
	jPk := JavaComPK{
		Package:    pkg,
		PKName:     pkName,
		Array:      arr,
		Capitalize: make([]string, len(arr)),
		AllGet:     "",
	}
	copy(jPk.Capitalize, jPk.Array)
	for i := range jPk.Capitalize {
		jPk.Capitalize[i] = strings.Title(jPk.Capitalize[i])
		if i == len(jPk.Capitalize)-1 {
			jPk.AllGet += "get" + strings.Title(jPk.Capitalize[i]) + "()"
		} else {
			jPk.AllGet += "get" + strings.Title(jPk.Capitalize[i]) + "(),"
		}
	}
	return &jPk
}

type JavaTemplate struct {
	Package   string
	TableName string
	Array     []string
	IDFields  []string
	IDClass   string
	TableRef  []TableRef
}

type TableRef struct {
	Name        string
	JoinColumns []ColumnRef
}

type ColumnRef struct {
	Col           string
	ReferencedCol string
}

func NewJavaTemplate(pkg, table string, idField, arr []string) *JavaTemplate {
	jTem := JavaTemplate{
		Package:   pkg,
		TableName: table,
		Array:     arr,
		IDFields:  idField,
	}
	if len(idField) > 1 {
		jTem.IDClass = "@IdClass(" + strings.Title(jTem.TableName) + "PK.class)"
	}
	return &jTem
}
