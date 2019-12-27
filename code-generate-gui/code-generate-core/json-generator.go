//USAGE INSTRUCTIONS
//Example: "go run main.go input.json rootPath", without filename.json provided, default filename will be "input.json"...
package code_generate_core

import (
	"archive/zip"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const (
	defaultFileName       = "input.json"
	defaultRootPath       = ""
	defaultProjectName    = "evaluation"
	defaultTemplateFolder = "code-generate-core/template"
)

type Input struct {
	Folders []Folder `json:"folders"`
}

type Folder struct {
	Env    []string
	Entity []string    `json:"entity"`
	RawEnv []string    `json:"env"`
	Model  string      `json:"model"`
	Files  []ModelJSON `json:"files"`
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

var input Input
var output Output
var templateDir string
var projectName string

func InputFileToInputStruct(filename string) string {
	jsonFile, err := os.Open(filename)
	if err != nil {
		return err.Error()
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err.Error()
	}
	err = jsonFile.Close()
	if err != nil {
		return err.Error()
	}

	err = json.Unmarshal(byteValue, &input)
	if err != nil {
		return err.Error()
	}
	return ""
}
func InputStringToInputStruct(guiInput string) string {

	err := json.Unmarshal([]byte(guiInput), &input)
	if err != nil {
		return err.Error()
	}
	return ""
}
func InputStructToOutputString(result *string) string {
	//WRITE THE OUT STRUCT
	output.RootPath = defaultRootPath
	output.ProjectName = projectName
	output.RootPath = strings.TrimSuffix(output.RootPath, "/")
	for k := range input.Folders {
		for i := range input.Folders[k].RawEnv {
			//Convert RawEnv to Env
			tmp := strings.LastIndex(input.Folders[k].RawEnv[i], "/")
			input.Folders[k].Env = append(input.Folders[k].Env, input.Folders[k].RawEnv[i][tmp+1:])

			//READ THE TEMPLATE FILES
			content, err := ioutil.ReadFile(defaultTemplateFolder + "/" + input.Folders[k].Env[i] + ".txt")
			if err != nil {
				return err.Error()
			}
			template := string(content)
			if strings.Contains(template, "{begin}") {
				text := template
				text = strings.ReplaceAll(text, "{env}", input.Folders[k].Env[i])
				for strings.Contains(text, "{begin}") {
					begin := strings.Index(text, "{begin}")
					end := strings.Index(text, "{end}")
					//envCount := strings.Count(text[begin:end], "{env}")
					entityCount := strings.Count(text[begin:end], "{entity}")
					entityLowerFirstCharacterCount := strings.Count(text[begin:end], "{entityLowerFirstCharacter}")
					tmpText := text[:end+len("{end}")]
					for j := 0; j < len(input.Folders[k].Entity)-1; j++ {
						tmpText += text[begin+len("{begin}") : end-1]
					}
					text = tmpText + text[end+len("{end}"):]

					for j := range input.Folders[k].Entity {
						//text = strings.Replace(text, "{env}", input.Folders[k].Env[i], envCount)
						text = strings.Replace(text, "{entity}", input.Folders[k].Entity[j], entityCount)
						text = strings.Replace(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:], entityLowerFirstCharacterCount)
					}
					text = strings.Replace(text, "{begin}", "", 1)
					text = strings.Replace(text, "{end}", "", 1)
				}
				filename := FileNameConverter(strings.ToUpper(output.ProjectName[:1])+output.ProjectName[1:], input.Folders[k].RawEnv[i]+"s")
				output.Files = append(output.Files, File{strings.ReplaceAll(input.Folders[k].RawEnv[i], "_", "-") + "/" + filename, text})
			} else {
				for j := range input.Folders[k].Entity {
					text := template
					text = strings.ReplaceAll(text, "{env}", input.Folders[k].Env[i])
					text = strings.ReplaceAll(text, "{entity}", input.Folders[k].Entity[j])
					text = strings.ReplaceAll(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:])
					filename := FileNameConverter(input.Folders[k].Entity[j], input.Folders[k].RawEnv[i])
					output.Files = append(output.Files, File{strings.ReplaceAll(input.Folders[k].RawEnv[i], "_", "-") + "/" + filename, text})
				}
			}
		}
		ModelJSONFileGenerator(FilesDetails{
			Model: input.Folders[k].Model,
			Files: input.Folders[k].Files,
		}, &output)
	}
	//OUTPUT STRUCT TO STRING(JSON)
	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return err.Error()
	}
	*result = string(file)
	return ""
}
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
func OutputStructToFiles(direc string) string {
	//CREATE FOLDER ON DISK
	if len(output.Files) == 0 {
		return "0 File Created On Disk"
	}
	if direc != "" {
		output.RootPath = direc
	}
	if output.RootPath != "" {
		err := os.MkdirAll(output.RootPath, os.ModePerm)
		if err != nil {
			return err.Error()
		}
		output.RootPath += "/"
	}
	for i := range output.Files {
		tmpPath := output.RootPath + output.Files[i].Name
		tmp := strings.LastIndex(tmpPath, "/")
		tmpPath = tmpPath[:tmp]
		err := os.MkdirAll(tmpPath, os.ModePerm)
		if err != nil {
			return err.Error()
		}
		f, err := os.Create(output.RootPath + output.Files[i].Name)
		if err != nil {
			return err.Error()
		}
		_, err = f.WriteString(output.Files[i].Content)
		if err != nil {
			return err.Error()
		}
		err = f.Close()
		if err != nil {
			return err.Error()
		}
	}
	//return strconv.Itoa(len(output.Files)) + " files created on disk"
	return ""
}
func OutputStructToZip(direc string) string {
	if len(output.Files) == 0 {
		return "No File To Zip"
	}
	fileName := output.ProjectName
	if direc != "" {
		tmp := strings.LastIndex(direc, "/")
		fileName = direc[tmp+1:]
		fileName = strings.TrimSuffix(fileName, ".zip")
		if tmp != -1 {
			output.RootPath = direc[:tmp]
		}
	}
	if output.RootPath != "" {
		err := os.MkdirAll(output.RootPath, os.ModePerm)
		if err != nil {
			return err.Error()
		}
		output.RootPath += "/"
	}

	newZipFile, err := os.Create(output.RootPath + fileName + ".zip")
	if err != nil {
		return err.Error()
	}
	defer newZipFile.Close()
	w := zip.NewWriter(newZipFile)
	for i := range output.Files {
		output.Files[i].Name = strings.TrimPrefix(output.Files[i].Name, "/")
		f, err := w.Create(output.Files[i].Name)
		if err != nil {
			return err.Error()
		}
		_, err = f.Write([]byte(output.Files[i].Content))
		if err != nil {
			return err.Error()
		}
	}
	err = w.Close()
	if err != nil {
		return err.Error()
	}
	//return "Zip created on disk"
	return ""
}
func GenerateFromString(temp, project, guiInput string, outputString *string) string {
	input = Input{}
	output = Output{}
	if temp != "" {
		templateDir = temp
	} else {
		templateDir = defaultTemplateFolder
	}
	if project != "" {
		projectName = project
	} else {
		projectName = defaultProjectName
	}
	err := InputStringToInputStruct(guiInput)
	if err != "" {
		return err
	}
	err = InputStructToOutputString(outputString)
	if err != "" {
		return err
	}
	return ""
}
func GenerateFromFile(temp, project, filename string, outputString *string) string {
	input = Input{}
	output = Output{}
	if temp != "" {
		templateDir = temp
	} else {
		templateDir = defaultTemplateFolder
	}
	if project != "" {
		projectName = project
	} else {
		projectName = defaultProjectName
	}
	err := InputFileToInputStruct(filename)
	if err != "" {
		return err
	}
	err = InputStructToOutputString(outputString)
	if err != "" {
		return err
	}
	return ""
}