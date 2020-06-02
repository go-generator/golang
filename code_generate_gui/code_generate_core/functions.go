package code_generate_core

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	. "github.com/go-generator/metadata"
)

func ToLower(s string) string {
	if s == "" {
		return ""
	}
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func AddStructFieldName(name string) string {
	return strings.Title(name)
}

func AddJSONTag(name string) string {
	return "`json:\"" + ToLower(name) + "\""
}

func AddBSONTag(name string) string {
	if strings.ToUpper(name) == "ID" {
		return " bson:\"" + "_id" + "\""
	}
	return " bson:\"" + ToLower(name) + "\""
}

func AddGORMTag(name string) string {
	return " gorm:\"column:" + ToLower(name) + "\"`\n"
}

func AddGORMPrimaryTag(name string) string {
	return " gorm:\"column:" + ToLower(name) + ":primary_key\"`\n"
}

func StandardizeStructName(s string) string {
	var res strings.Builder
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	tokens := strings.Split(s, "_")
	for _, v := range tokens {
		alphanumericString := reg.ReplaceAllString(v, "")
		res.WriteString(strings.Title(alphanumericString))
	}
	return res.String()
}

func WritePackage(m *Model, packageName string) {
	m.WriteFile.WriteString("package " + packageName + "\n\n")
}

func WriteTypeAlias(m *Model) {
	for _, v := range m.Alias {
		m.WriteFile.WriteString("type " + v.Name + " " + v.Type + "\n\n")
	}
}

func WriteStruct(m *Model) FileInfo {
	var count int
	tmpl := FileInfo{}
	tmpl.Name = m.Name
	tmpl.StructName = StandardizeStructName(m.Name)
	for _, v := range m.Fields {
		tmp := FieldInfo{
			Name: StandardizeStructName(v.Name),
			Type: v.Type,
		}
		if v.Id {
			count++
			tmpl.IDFields = append(tmpl.IDFields, tmp)
		} else {
			tmpl.Fields = append(tmpl.Fields, tmp)
		}
	}
	return tmpl
}
