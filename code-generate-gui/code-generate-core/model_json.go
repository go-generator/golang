package code_generate_core

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type FilesDetails struct {
	Model string      `json:"model"`
	Files []ModelJSON `json:"files"`
}

type ModelJSON struct {
	Env        string          `json:"env"`
	Name       string          `json:"name"`
	Source     string          `json:"source"`
	ConstValue []Const         `json:"const"`
	TypeAlias  []TypeAlias     `json:"type_alias"`
	Fields     []FieldElements `json:"fields"`
	WriteFile  strings.Builder // Writing content of the file
}

func ToLower(s string) string {
	if s == "" {
		return ""
	}
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
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
	PrimaryKey bool   `json:"primaryKey"`
}

func ReadJSON(pathFile string) FilesDetails {
	var v FilesDetails
	jsonFile, err := os.Open(pathFile)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &v)
	if err != nil {
		panic(err)
	}
	return v
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

func (m *ModelJSON) WritePackage(packageName string) string {
	m.WriteFile.WriteString("package " + packageName + "\n\n")
	return "package " + packageName + "\n\n"
}

func (m *ModelJSON) WriteTypeAlias() {
	for _, v := range m.TypeAlias {
		m.WriteFile.WriteString("type " + v.Name + " " + v.Type + "\n\n")
	}
}

func (m *ModelJSON) WriteConstValue() {
	m.WriteFile.WriteString("const (\n")
	for _, v := range m.ConstValue {
		switch v.Value.(type) {
		case string:
			m.WriteFile.WriteString("\t" + v.Name + " " + v.Type + " = " + "\"" + v.Value.(string) + "\"" + "\n")
		default:
			m.WriteFile.WriteString("\t" + v.Name + " " + v.Type + " = " + fmt.Sprint(v.Value) + "\n")
		}
	}
	m.WriteFile.WriteString(")\n\n")
}

func (m *ModelJSON) WriteStruct() {
	var count int
	for _, v := range m.Fields {
		if v.PrimaryKey {
			count++
		}
	}
	m.WriteFile.WriteString("type " + m.Name + " struct {\n")
	if count < 2 {
		for _, v := range m.Fields {
			if v.PrimaryKey {
				m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag("_id") + AddGORMTag(v.Source, true))
				continue
			}
			m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag(v.Name) + AddGORMTag(v.Source, false))
		}
	} else {
		for _, v := range m.Fields {
			if v.PrimaryKey {
				m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag(v.Name) + AddGORMTag(v.Source, true))
				continue
			}
			m.WriteFile.WriteString("\t" + AddStructFieldName(v.Name) + "\t" + v.Type + "\t" + AddJSONTag(v.Name) + AddBSONTag(v.Name) + AddGORMTag(v.Source, false))
		}
	}
	m.WriteFile.WriteString("}")
}

func (m *ModelJSON) CreateContent(packageName string) {
	m.WritePackage(packageName)
	m.WriteTypeAlias()
	m.WriteConstValue()
	m.WriteStruct()
}

func ReadAllSubFiles(rootPath string) []string {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func ModelJSONFileGenerator(content FilesDetails, out *Output) {
	//var out Output
	//out.ProjectName = projectName
	//out.RootPath = rootPath
	//if _, err := os.Stat(destination); os.IsNotExist(err) {
	//	err = os.Mkdir(destination, 0777)
	//	if err != nil {
	//		log.Fatal("Failed attempt to create directory, " + err.Error())
	//	}
	//}
	//fileDirectory := destination + output + ".json"
	//jsonFiles := ReadAllSubFiles(source)

	var file File
	//content := ReadJSON(source + v) // struct content

	for _, k := range content.Files {
		k.CreateContent(content.Model)
		file.Name = content.Model + "/" + ToLower(k.Name) + ".go"
		file.Content = k.WriteFile.String()
		out.Files = append(out.Files, file)
	}

	//data, err := json.MarshalIndent(out, "", " ")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = ioutil.WriteFile(fileDirectory, data, 0644) // Create and write files
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("Generated Successfully")
}
