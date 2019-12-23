//USAGE INSTRUCTIONS
//Example: "go run main.go input.json template.txt rootPath", without filename.json provided, default filename will be "input.json"...
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const (
	defaultFileName    = "input.json"
	defaultTemplate    = "template.txt"
	defaultRootPath    = ""
	defaultProjectName = "project"
)

type Input struct {
	Env    []string `json:"env"`
	Entity []string `json:"entity"`
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

func main() {

	//READ THE JSON INPUT

	filename := defaultFileName
	templateFile := defaultTemplate
	rootPath := defaultRootPath
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	if len(os.Args) > 2 {
		templateFile = os.Args[2]
	}
	if len(os.Args) > 3 {
		rootPath = os.Args[3]
		rootPath = strings.TrimSuffix(rootPath, "/")
	}
	jsonFile, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	err = jsonFile.Close()
	if err != nil {
		panic(err)
	}
	var input Input
	var output Output
	err = json.Unmarshal(byteValue, &input)
	if err != nil {
		panic(err)
	}

	//READ THE TEMPLATE FILE
	content, err := ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}
	template := string(content)

	//WRITE THE JSON OUTPUT
	for i := range input.Env {
		for j := range input.Entity {
			text := template
			text = strings.ReplaceAll(text, "{env}", input.Env[i])
			text = strings.ReplaceAll(text, "{entity}", input.Entity[j])
			text = strings.ReplaceAll(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Entity[j])[0])+input.Entity[j][1:])
			filename := FileNameConverter(input.Entity[j], rootPath, input.Env[i])
			output.Files = append(output.Files, File{input.Env[i] + "/" + filename, text})
		}
	}
	output.RootPath = rootPath
	output.ProjectName = defaultProjectName
	file, _ := json.MarshalIndent(output, "", " ")
	_ = ioutil.WriteFile("output.json", file, 0644)
}

//Convert SomeThing to some_thing
func FileNameConverter(s string, rootPath string, ss string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "_" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	rootPath = strings.ReplaceAll(rootPath, "/", "_")
	if rootPath != "" {
		return s3[1:] + "_" + rootPath + "_" + ss + ".go"
	} else {
		return s3[1:] + "_" + ss + ".go"
	}
}
