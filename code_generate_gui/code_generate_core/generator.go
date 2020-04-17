//USAGE INSTRUCTIONS
//Example: "go run main.go input.json rootPath", without filename.json provided, default filename will be "input.json"...
package code_generate_core

import (
	"archive/zip"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sqweek/dialog"
	"golang/code_generate_gui/code_generate_core/common"
	"golang/code_generate_gui/code_generate_core/model"
)

const (
	defaultFileName                   = "input.json"
	defaultRootPath                   = ""
	defaultProjectName                = "evaluation"
	defaultTemplateFolderRelativePath = "./code_generate_core/template"
)

var defaultTemplateFolder = absTemplatePath()

func absTemplatePath() string {
	absPath, err := filepath.Abs(defaultTemplateFolderRelativePath)
	if err != nil {
		log.Println(err)
	}
	return absPath
}

var (
	input       model.Input
	output      model.Output
	templateDir string
	projectName string
)

func InputJsonFileToInputStruct(filename string) string {
	byteValue, err := ioutil.ReadFile(filename)
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
			//Convert RawEnv to Model
			tmp := strings.LastIndex(input.Folders[k].RawEnv[i], "/")
			input.Folders[k].Env = append(input.Folders[k].Env, input.Folders[k].RawEnv[i][tmp+1:])

			//READ THE TEMPLATE FILES
			content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + input.Folders[k].Env[i] + ".txt")
			if err != nil {
				return err.Error()
			}
			template := string(content)
			if strings.Contains(template, "{begin}") {
				text := template
				text = strings.ReplaceAll(text, "{env}", input.Folders[k].Env[i])
				text = strings.ReplaceAll(text, "{projectName}", projectName)
				text = strings.ReplaceAll(text, "{projectNameUpperFirstCharacter}", strings.Title(projectName))
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
						//text = strings.Replace(text, "{env}", input.ModelFile[k].Model[i], envCount)
						text = strings.Replace(text, "{entity}", input.Folders[k].Entity[j], entityCount)
						text = strings.Replace(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:], entityLowerFirstCharacterCount)
					}
					text = strings.Replace(text, "{begin}", "", 1)
					text = strings.Replace(text, "{end}", "", 1)
				}
				filename := FileNameConverter(strings.ToUpper(output.ProjectName[:1])+output.ProjectName[1:], input.Folders[k].RawEnv[i]+"s")
				output.Files = append(output.Files, model.File{strings.ReplaceAll(input.Folders[k].RawEnv[i], "_", "-") + "/" + filename, text})
			} else {
				for j := range input.Folders[k].Entity {
					text := template
					text = strings.ReplaceAll(text, "{env}", input.Folders[k].Env[i])
					text = strings.ReplaceAll(text, "{projectName}", projectName)
					text = strings.ReplaceAll(text, "{projectNameUpperFirstCharacter}", strings.Title(projectName))
					text = strings.ReplaceAll(text, "{entity}", input.Folders[k].Entity[j])
					text = strings.ReplaceAll(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:])
					filename := FileNameConverter(input.Folders[k].Entity[j], input.Folders[k].RawEnv[i])
					output.Files = append(output.Files, model.File{strings.ReplaceAll(input.Folders[k].RawEnv[i], "_", "-") + "/" + filename, text})
				}
			}
		}
		FileDetailsToOutput(model.FilesDetails{
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

func OutputStructToFiles(directory string) string {
	//CREATE FOLDER ON DISK
	if len(output.Files) == 0 {
		return "0 File Created On Disk"
	}
	if directory != "" {
		output.RootPath = directory
	}
	if output.RootPath != "" {
		err := os.MkdirAll(output.RootPath, os.ModePerm)
		if err != nil {
			return err.Error()
		}
		output.RootPath += string(os.PathSeparator)
	}
	var allFiles []string
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
		allFiles = append(allFiles, f.Name())
		err = f.Close()
		if err != nil {
			return err.Error()
		}
	}
	var wg sync.WaitGroup
	for i := range allFiles {
		wg.Add(1)
		go func(dir string, wtg *sync.WaitGroup) {
			defer wtg.Done()
			_, err := ShellExecutor("goimports", []string{"-w", dir})
			if err != nil {
				log.Println(err)
			}
		}(allFiles[i], &wg)
	}
	wg.Wait()
	return ""
}

func OutputStructToZip() string {
	if output.Files == nil {
		file, err := dialog.File().Filter("json file", "json").Load()
		res := ""
		if err != nil {
			return err.Error()
		}
		GenerateFromFile(templateDir, projectName, file, &res)
	}
	directory, err := dialog.File().Filter("zip file", "zip").Title("Export to zip").Save()
	if err != nil {
		return err.Error()
	}
	fileName := output.ProjectName
	if directory != "" {
		tmp := strings.LastIndex(directory, string(os.PathSeparator))
		fileName = directory[tmp+1:]
		fileName = strings.TrimSuffix(fileName, ".zip")
		if tmp != -1 {
			output.RootPath = directory[:tmp]
		}
	}
	if output.RootPath != "" {
		err := os.MkdirAll(output.RootPath, os.ModePerm)
		if err != nil {
			return err.Error()
		}
		output.RootPath += string(os.PathSeparator)
	}
	tmp := filepath.Join([]string{".", "tmp.go"}...)
	for i := range output.Files {
		err := ioutil.WriteFile(tmp, []byte(output.Files[i].Content), 0664)
		if err != nil {
			return err.Error()
		}
		_, err = ShellExecutor("goimports", []string{"-w", tmp})
		if err != nil {
			return err.Error()
		}
		formattedData, err := ioutil.ReadFile(tmp)
		if err != nil {
			return err.Error()
		}
		output.Files[i].Content = string(formattedData)
	}
	err = os.Remove(tmp)
	if err != nil {
		return err.Error()
	}
	newZipFile, err := os.Create(output.RootPath + fileName + ".zip")
	if err != nil {
		return err.Error()
	}
	w := zip.NewWriter(newZipFile)
	defer func() string {
		err = newZipFile.Close()
		if err != nil {
			return err.Error()
		}
		return ""
	}()
	defer func() string {
		err = w.Close()
		if err != nil {
			return err.Error()
		}
		return ""
	}()
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
	return ""
}

func GenerateFromString(temp, project, guiInput string, outputString *string) string {
	input = model.Input{}
	output = model.Output{}
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
	input = model.Input{}
	output = model.Output{}
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
	err := InputJsonFileToInputStruct(filename)
	if err != "" {
		return err
	}
	err = InputStructToOutputString(outputString)
	if err != "" {
		return err
	}
	return ""
}

func ShellExecutor(program string, arguments []string) ([]byte, error) {
	cmd := exec.Command(program, arguments...)
	return cmd.Output()
}

func FileDetailsToOutput(content model.FilesDetails, out *model.Output) {
	var file model.File
	for _, k := range content.Files {
		common.CreateContent(&k, content.Model)
		file.Name = content.Model + "/" + common.ToLower(k.Name) + ".go"
		file.Content = k.WriteFile.String()
		out.Files = append(out.Files, file)
	}
}
