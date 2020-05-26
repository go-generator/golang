//USAGE INSTRUCTIONS
//Example: "go run main.go input.json rootPath", without filename.json provided, default filename will be "input.json"...
package code_generate_core

import (
	"archive/zip"
	"context"
	"encoding/json"
	"github.com/go-generator/generator"
	iou "github.com/go-generator/io"
	"github.com/go-generator/metadata"
	"github.com/go-generator/project"
	"github.com/sqweek/dialog"
	"golang/code_generate_gui/code_generate_core/model"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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
	input       metadata.Project
	output      model.Output
	templateDir string
	projectName string
)

func ShareMapInitF(packageName, projectName string) map[string]string {
	return map[string]string{
		"package_name": packageName,
		"projectName":  projectName,
		"ProjectName":  strings.Title(projectName),
	}
}

func FullMapInitF(env string, entity string) map[string]string {
	return map[string]string{
		"static_package": env,
		"projectName":    projectName,
		"ProjectName":    strings.Title(projectName),
		"Entity":         entity,
		"entity":         strings.ToLower(string(entity[0])) + entity[1:],
	}
}
func EnvTemplateF(template string, fullMap map[string]string) string {
	text := template
	for k, v := range fullMap {
		text = strings.ReplaceAll(text, k, v)
	}
	return text
}
func ArrayTemplateF(template string, share map[string]string, arr []map[string]string) string {
	text := template
	for k, v := range share {
		text = strings.ReplaceAll(text, k, v)
	}
	for strings.Contains(text, "{begin}") {
		begin := strings.Index(text, "{begin}")
		end := strings.Index(text, "{end}")
		tmpText := text[:begin]
		for j := 0; j < len(arr); j++ {
			tmp := text[begin+len("{begin}") : end-1]
			for k, v := range arr[j] {
				tmp = strings.ReplaceAll(tmp, k, v)
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
	input.Env = project.EnvInit(input.Env, "hotelManagement")
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
	output.ProjectName = input.Env["projectName"]
	output.RootPath = strings.TrimSuffix(output.RootPath, "/")
	templateMap, err := generator.TemplateMapInitF(defaultTemplateFolder, input)
	if err != nil {
		return err.Error()
	}

	var t generator.Generator
	t = &generator.JavaGenerator{}
	outputFiles := t.Generate(context.Background(), input, templateMap)
	output.Files = outputFiles
	//OUTPUT STRUCT TO STRING(JSON)
	file, err := json.MarshalIndent(output, "", " ")
	if err != nil {
		return err.Error()
	}
	*result = string(file)
	return ""
}

func FileNameConverter(s string, path string) string {
	ext := ".go"
	if strings.Contains(path, ".") {
		ext = ""
	}
	path = strings.ReplaceAll(path, "/", "_")
	if s == "" {
		return path + ext
	}
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if s2[i] != s[i] {
			s3 += "_" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}

	return s3[1:] + "_" + path + ext
}
func ConvertFileStructFromMetadataToIoFormat(f metadata.File) iou.File {
	return iou.File{
		Name:    f.Name,
		Content: f.Content,
	}
}
func ConvertListFileStructFromMetadataToIoFormat(f []metadata.File) []iou.File {
	var out []iou.File
	for _, v := range f {
		out = append(out, ConvertFileStructFromMetadataToIoFormat(v))
	}
	return out
}
func OutputStructToFiles(directory string) string {
	err := iou.SaveFiles(directory, ConvertListFileStructFromMetadataToIoFormat(output.Files))
	if err != nil {
		return err.Error()
	}
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
	err = iou.SaveFilesToZip(output.RootPath+fileName+".zip", ConvertListFileStructFromMetadataToIoFormat(output.Files))
	//err = iou.SaveFileToZip(output.RootPath+fileName+".zip", output.Files[0].Name,output.Files[0].Content)
	if err != nil {
		return err.Error()
	}
	return ""
}
func OldOutputStructToZip() string {
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

//func OldGenerateFromString(temp, project, guiInput string, outputString *string) string {
//	input =
//	output = model.Output{}
//	if temp != "" {
//		templateDir = temp
//	} else {
//		templateDir = defaultTemplateFolder
//	}
//	if project != "" {
//		projectName = project
//	} else {
//		projectName = defaultProjectName
//	}
//	err := InputStringToInputStruct(guiInput)
//	if err != "" {
//		return err
//	}
//	err = InputStructToOutputString(outputString)
//	if err != "" {
//		return err
//	}
//	return ""
//}
func GenerateFromString(temp, project, guiInput string, outputString *string) string {
	return ""
}
func GenerateFromFile(temp, project, filename string, outputString *string) string {
	input = metadata.Project{}
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
	var file metadata.File
	for _, k := range content.Files {
		out.OutFile = append(output.OutFile, WriteStruct(&k))
		file.Name = content.Model + "/" + ToLower(k.Name) + ".go"
		file.Content = k.WriteFile.String()
		out.Files = append(out.Files, file)
	}
}
