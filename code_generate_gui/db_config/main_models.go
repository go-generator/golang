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
	Name          string          `json:"name"`
	Source        string          `json:"source"`
	ConstValue    []Const         `json:"const"`
	TypeAlias     []TypeAlias     `json:"type_alias"`
	Relationships []Relationship  `json:"arrays"`
	Fields        []FieldElements `json:"fields"`
	WriteFile     strings.Builder `json:"-"`
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
	Name       string `json:"name"`
	Source     string `json:"source"`
	Type       string `json:"type"`
	PrimaryKey bool   `json:"primaryKey,omitempty"`
}

type Relationship struct {
	Table  string  `json:"table"`
	Fields []Field `json:"fields"`
}

func (m *ModelJSON) ExtractFieldType() map[string]string {
	res := make(map[string]string)
	for _, v := range m.Fields {
		res[v.Name] = v.Type
	}
	return res
}

type Connection struct {
	TableName       string
	ReferencedTable string
	Fields          []Field
}

type Field struct {
	ColumnName       string `json:"column"`
	ReferencedColumn string `json:"ref"`
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
