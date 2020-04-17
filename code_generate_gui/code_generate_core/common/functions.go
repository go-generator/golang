package common

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
	return " bson:\"" + ToLower(name) + "\""
}

func AddGORMTag(name string, primaryTag bool) string {
	if name == "" {
		return "`\n"
	}
	if primaryTag {
		return " gorm:\"column:" + name + ":primary_key\"`\n"
	}
	return " gorm:\"column:" + name + "\"`\n"
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

func WriteStruct(m *Model) {
	var count int
	for _, v := range m.Fields {
		if v.Id {
			count++
		}
	}
	m.WriteFile.WriteString("type " + StandardizeStructName(m.Name) + " struct {\n")
	if count < 2 {
		for _, v := range m.Fields {
			if v.Id {
				m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag("_id") + AddGORMTag(v.Source, true))
				continue
			}
			m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag(v.Name) + AddGORMTag(v.Source, false))
		}
	} else {
		for _, v := range m.Fields {
			if v.Id {
				m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag(v.Name) + AddGORMTag(v.Source, true))
				continue
			}
			m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag(v.Name) + AddGORMTag(v.Source, false))
		}
	}
	m.WriteFile.WriteString("}")
}

func CreateContent(m *Model, packageName string) {
	WritePackage(m, packageName)
	WriteTypeAlias(m)
	WriteStruct(m)
}
