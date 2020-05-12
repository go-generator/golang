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
	"text/template"

	iou "github.com/go-generator/io"
	"github.com/sqweek/dialog"
	"golang/code_generate_gui/code_generate_core/model"
)

const (
	defaultFileName                   = "input.json"
	defaultRootPath                   = ""
	defaultProjectName                = "evaluation"
	defaultTemplateFolderRelativePath = "./template"
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
func ShareMapInitF(env string) map[string]string{
	return map[string]string{
		"${env}": env,
		"${projectName}": projectName,
		"${projectNameUpperFirstCharacter}": strings.Title(projectName),
	}
}
func ArrMapInitF(entityList []string) []map[string]string{
	var tmp []map[string]string
	for _,v:= range entityList {
		tmp = append(tmp,
			map[string]string{
				"${entity}":                    v,
				"${entityLowerFirstCharacter}": strings.ToLower(string(v[0])) + v[1:],
			})
	}
	return tmp
}
func FullMapInitF(env string, entity string) map[string]string{
	return map[string]string{
		"${env}": env,
		"${projectName}": projectName,
		"${projectNameUpperFirstCharacter}": strings.Title(projectName),
		"${entity}":                    entity,
		"${entityLowerFirstCharacter}": strings.ToLower(string(entity[0])) + entity[1:],
	}
}
func EnvTemplateF(template string, fullMap map[string]string) string {
	text := template
	for k,v:=range fullMap{
		text=strings.ReplaceAll(text,k,v)
	}
	return text
}
func ArrayTemplateF(template string, share map[string]string, arr []map[string]string ) string{
	text := template
	for k,v:=range share{
		text = strings.ReplaceAll(text, k, v)
	}
	for strings.Contains(text, "{begin}") {
		begin := strings.Index(text, "{begin}")
		end := strings.Index(text, "{end}")
		tmpText := text[:begin]
		for j := 0; j < len(arr); j++ {
			tmp:= text[begin+len("{begin}") : end-1]
			for k,v:=range arr[j]{
				tmp=strings.ReplaceAll(tmp,k,v)
			}
			tmpText += tmp
		}
		text = tmpText + text[end+len("{end}"):]
		//text = strings.Replace(text, "{begin}", "", 1)
		//text = strings.Replace(text, "{end}", "", 1)
	}
	return text
}

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
			CleanEnv:=input.Folders[k].RawEnv[i][tmp+1:]

			//READ THE TEMPLATE FILES
			content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + CleanEnv + ".txt")
			if err != nil {
				return err.Error()
			}
			template := string(content)
			for j := range input.Folders[k].Entity {
				text := EnvTemplateF(template, FullMapInitF(CleanEnv,input.Folders[k].Entity[j]))
				filename := FileNameConverter(input.Folders[k].Entity[j], input.Folders[k].RawEnv[i])
				output.Files = append(output.Files, iou.File{strings.ReplaceAll(input.Folders[k].RawEnv[i], "_", "-") + "/" + filename, text})
			}
			}
		for i := range input.Folders[k].Array {
			//Convert RawEnv to Model
			tmp := strings.LastIndex(input.Folders[k].Array[i], "/")
			CleanEnv:=input.Folders[k].Array[i][tmp+1:]

			//READ THE TEMPLATE FILES
			content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + CleanEnv + ".txt")
			if err != nil {
				return err.Error()
			}
			template := string(content)
			text:= ArrayTemplateF(template, ShareMapInitF(CleanEnv), ArrMapInitF(input.Folders[k].Entity))
			filename := FileNameConverter(strings.ToUpper(output.ProjectName[:1])+output.ProjectName[1:], input.Folders[k].Array[i]+"s")
			output.Files = append(output.Files, iou.File{strings.ReplaceAll(input.Folders[k].Array[i], "_", "-") + "/" + filename, text})

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
	fl := output
	log.Println(fl)
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
	modelTmpl := filepath.Join(defaultTemplateFolder, "model.txt")
	funcMap := template.FuncMap{
		"AddJsonTag":        AddJSONTag,
		"AddBsonTag":        AddBSONTag,
		"AddGormTag":        AddGORMTag,
		"AddGormPrimaryTag": AddGORMPrimaryTag,
	}
	tmplName := filepath.Base(modelTmpl)
	mTmpl := template.New(tmplName).Funcs(funcMap)
	mT, err := mTmpl.ParseFiles(modelTmpl)
	if err != nil {
		return err.Error()
	}
	for i := range output.OutFile {
		err := os.MkdirAll(filepath.Join(output.RootPath, "model"), os.ModePerm)
		if err != nil {
			return err.Error()
		}
		f, err := os.Create(filepath.Join(output.RootPath, "model", output.OutFile[i].Name+".go"))
		if err != nil {
			return err.Error()
		}
		err = mT.Execute(f, output.OutFile[i])
		if err != nil {
			return err.Error()
		}
		allFiles = append(allFiles, f.Name())
	}
	for i := range output.Files {
		if output.Files[i].Content != "" {
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
	var file iou.File
	for _, k := range content.Files {
		out.OutFile = append(output.OutFile, WriteStruct(&k))
		file.Name = content.Model + "/" + ToLower(k.Name) + ".go"
		file.Content = k.WriteFile.String()
		out.Files = append(out.Files, file)
	}
}
