package model_json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	. "github.com/go-generator/metadata"
)

type FilesDetails struct {
	Env   string  `json:"env"`
	Files []Model `json:"files"`
}

func ToLower(s string) string {
	if s == "" {
		return ""
	}
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
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

func WritePackage(m *Model, packageName string) {
	m.WriteFile.WriteString("package " + packageName + "\n\n")
	for _, v := range m.Fields {
		if v.Type == "time.Time" {
			m.WriteFile.WriteString("import \"time\"\n\n")
			break
		}
	}
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

type Output struct {
	ProjectName string `json:"projectName"`
	RootPath    string `json:"rootPath"`
	Files       []File `json:"files"`
}

type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func StandardizeStructName(s string) string {
	var res strings.Builder
	tokens := strings.Split(s, "_")
	for _, v := range tokens {
		res.WriteString(strings.Title(v))
	}
	return res.String()
}

func ModelJSONFileGenerator(source, destination, projectName, rootPath, output string) {
	var out Output
	out.ProjectName = projectName
	out.RootPath = rootPath
	if _, err := os.Stat(destination); os.IsNotExist(err) {
		err = os.Mkdir(destination, 0777)
		if err != nil {
			log.Fatal("Failed attempt to create directory, " + err.Error())
		}
	}
	fileDirectory := destination + output + ".json"
	jsonFiles := ReadAllSubFiles(source)
	for _, v := range jsonFiles {
		var file File
		content := ReadJSON(source + v) // struct content
		for _, k := range content.Files {
			CreateContent(&k, content.Env)
			file.Name = content.Env + "/" + ToLower(k.Name) + ".go"
			file.Content = k.WriteFile.String()
			out.Files = append(out.Files, file)
		}
	}
	data, err := json.MarshalIndent(out, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(fileDirectory, data, 0644) // Create and write files
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Generated Successfully")
}
