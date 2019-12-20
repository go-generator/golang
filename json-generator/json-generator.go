//USAGE INSTRUCTIONS
//Example: "go run main.go filename.json", without filename.json provided, default filename will be "input.json"
package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const (
	defaultFileName = "input.json"
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
	if len(os.Args) > 1 {
		filename = os.Args[1]
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
	content, err := ioutil.ReadFile("template.txt")
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
			filename := FileNameConverter(input.Entity[j], input.Env[i])
			output.Files = append(output.Files, File{input.Env[i] + "/" + filename, text})
		}
	}
	file, _ := json.MarshalIndent(output, "", " ")
	_ = ioutil.WriteFile("output.json", file, 0644)
}
func FileNameConverter(s string, ss string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "_" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	return s3[1:] + "_" + ss + ".go"
}
