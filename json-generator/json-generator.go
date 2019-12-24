//USAGE INSTRUCTIONS
//Example: "go run main.go input.json rootPath", without filename.json provided, default filename will be "input.json"...
package main

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	defaultFileName       = "input.json"
	defaultRootPath       = ""
	defaultProjectName    = "evaluation"
	defaultTemplateFolder = "template"
)

type Input struct {
	Folders []Folder `json:"folders"`
}

type Folder struct {
	Env    []string
	Entity []string `json:"entity"`
	RawEnv []string `json:"env"`
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
	rootPath := defaultRootPath
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	if len(os.Args) > 2 {
		rootPath = os.Args[2]
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

	//WRITE THE OUT STRUCT
	output.RootPath = rootPath
	output.ProjectName = defaultProjectName
	for k := range input.Folders {
		for i := range input.Folders[k].RawEnv {
			//Convert RawEnv to Env
			tmp := strings.LastIndex(input.Folders[k].RawEnv[i], "/")
			input.Folders[k].Env = append(input.Folders[k].Env, input.Folders[k].RawEnv[i][tmp+1:])

			//READ THE TEMPLATE FILES
			content, err := ioutil.ReadFile(defaultTemplateFolder + "/" + input.Folders[k].Env[i] + ".txt")
			if err != nil {
				panic(err)
			}
			template := string(content)
			if strings.Contains(template, "{begin}") {
				text := template
				for strings.Contains(text, "{begin}") {
					begin := strings.Index(text, "{begin}")
					end := strings.Index(text, "{end}")
					envCount := strings.Count(text[begin:end], "{env}")
					entityCount := strings.Count(text[begin:end], "{entity}")
					entityLowerFirstCharacterCount := strings.Count(text[begin:end], "{entityLowerFirstCharacter}")
					tmpText := text[:end+len("{end}")]
					for j := 0; j < len(input.Folders[k].Entity)-1; j++ {
						tmpText += text[begin+len("{begin}") : end-1]
					}
					text = tmpText + text[end+len("{end}"):]

					for j := range input.Folders[k].Entity {
						text = strings.Replace(text, "{env}", input.Folders[k].Env[i], envCount)
						text = strings.Replace(text, "{entity}", input.Folders[k].Entity[j], entityCount)
						text = strings.Replace(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:], entityLowerFirstCharacterCount)
					}
					text = strings.Replace(text, "{begin}", "", 1)
					text = strings.Replace(text, "{end}", "", 1)
				}
				filename := FileNameConverter(strings.ToUpper(output.ProjectName[:1])+output.ProjectName[1:], input.Folders[k].RawEnv[i]+"s")
				output.Files = append(output.Files, File{input.Folders[k].RawEnv[i] + "/" + filename, text})
			} else {
				for j := range input.Folders[k].Entity {
					text := template
					text = strings.ReplaceAll(text, "{env}", input.Folders[k].Env[i])
					text = strings.ReplaceAll(text, "{entity}", input.Folders[k].Entity[j])
					text = strings.ReplaceAll(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:])
					filename := FileNameConverter(input.Folders[k].Entity[j], input.Folders[k].RawEnv[i])
					output.Files = append(output.Files, File{input.Folders[k].RawEnv[i] + "/" + filename, text})
				}
			}
		}
	}

	output.RootPath = strings.TrimSuffix(output.RootPath, "/")
	if output.RootPath != "" {
		err = os.MkdirAll(output.RootPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		output.RootPath += "/"
	}

	//CREATE FOLDER ON DISK
	for i := range output.Files {
		tmpPath := output.RootPath + output.Files[i].Name
		tmp := strings.LastIndex(tmpPath, "/")
		tmpPath = tmpPath[:tmp]
		err = os.MkdirAll(tmpPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(output.RootPath + output.Files[i].Name)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(output.Files[i].Content)
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(strconv.Itoa(len(output.Files)) + " files created on disk")

	//New-style ZIP
	if output.ProjectName == "" {
		output.ProjectName = defaultProjectName
	}
	newZipFile, err := os.Create(output.RootPath + output.ProjectName + ".zip")
	if err != nil {
		panic(err)
	}
	defer newZipFile.Close()
	w := zip.NewWriter(newZipFile)
	for i := range output.Files {
		output.Files[i].Name = strings.TrimPrefix(output.Files[i].Name, "/")
		f, err := w.Create(output.Files[i].Name)
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(output.Files[i].Content))
		if err != nil {
			panic(err)
		}
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Zip created on disk")

	////CREATE OUTPUT JSON
	//file, _ := json.MarshalIndent(output, "", " ")
	//_ = ioutil.WriteFile("output.json", file, 0644)
}

//Convert (SomeThing,service/impl) to some_thing_service_impl.go
func FileNameConverter(s string, path string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "_" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	path = strings.ReplaceAll(path, "/", "_")
	return s3[1:] + "_" + path + ".go"
}
