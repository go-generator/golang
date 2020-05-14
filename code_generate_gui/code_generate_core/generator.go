//USAGE INSTRUCTIONS
//Example: "go run main.go input.json rootPath", without filename.json provided, default filename will be "input.json"...
package code_generate_core

import (
	"archive/zip"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"text/template"

	"github.com/go-generator/metadata"
	templ "github.com/go-generator/template"
	"github.com/go-generator/generator"
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

func ShareMapInitF(packageName, projectName string) map[string]string {
	return map[string]string{
		"package_name": packageName,
		"projectName":  projectName,
		"ProjectName":  strings.Title(projectName),
	}
}
func ArrMapInitF(entityList []string) []map[string]string {
	var mapList []map[string]string
	for _, v := range entityList {
		mapList = append(mapList, templ.BuildNames(v))
	}
	return mapList
}
func ModelArrMapInitF(fieldList []metadata.Field) []map[string]string {
	var mapList []map[string]string
	for _, v := range fieldList {
		mapList = append(mapList, templ.BuildFields(v.Name, v.Id, 0, v.Type))
	}
	return mapList
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

type DefaultEntityTemplate struct {
}

func (t *DefaultEntityTemplate) Merge(ctx context.Context, template string, share map[string]string, parent map[string]string, fields []map[string]string) string {
	s := template
	for k, v := range share {
		s = strings.ReplaceAll(s, "${env:"+k+"}", v)
	}
	for k, v := range parent {
		s = strings.ReplaceAll(s, "${"+k+"}", v)
	}
	return s
}

type DefaultGenerator struct {
}

func (t *DefaultGenerator) Generate(ctx context.Context, project metadata.Project, templates map[string]string, types map[string]string) []metadata.File {
	var outputFile []metadata.File
	for _, v := range project.Statics {
		//COMPLETE THE FILENAME
		var t templ.DefaultStaticTemplate
		v.File = t.Generate(context.Background(), v.File, project.Env)
		//READ THE TEMPLATE FILES
		template := templates[v.Name]
		//CREATE TEXT
		text := t.Generate(context.Background(), template, project.Env)
		outputFile = append(outputFile, metadata.File{"/" + v.File, text})
	}
	for _, v := range project.Arrays {
		//COMPLETE THE FILENAME
		var t1 templ.DefaultStaticTemplate
		v.File = t1.Generate(context.Background(), v.File, project.Env)
		//READ THE TEMPLATE FILES
		template := templates[v.Name]
		//CREATE TEXT
		t := templ.NewArrayTemplate("${begin}", "${end}", "", "")
		text := t.Array(context.Background(), template, project.Env, nil, ArrMapInitF(project.Collection))
		outputFile = append(outputFile, metadata.File{"/" + v.File, text})
	}
	for _, v := range project.Entities {
		for _, c := range project.Collection {
			buildNamesMap := templ.BuildNames(c)
			//COMPLETE THE FILENAME
			var t templ.EntityTemplate
			t = &DefaultEntityTemplate{}
			tmpFile := t.Merge(context.Background(), v.File, project.Env, buildNamesMap, nil)

			//READ THE TEMPLATE FILES
			template := templates[v.Name]
			//CREATE TEXT
			text := t.Merge(context.Background(), template, project.Env, buildNamesMap, nil)
			outputFile = append(outputFile, metadata.File{"/" + tmpFile, text})
		}
	}
	for _, v := range project.Models {
		//READ THE TEMPLATE FILES
		template := templates["model"]
		//CREATE TEXT
		var primaryFields []metadata.Field
		var noPrimaryFields []metadata.Field
		for _, v1 := range v.Fields {
			if v1.Id {
				primaryFields = append(primaryFields, v1)
			} else {
				noPrimaryFields = append(noPrimaryFields, v1)
			}
		}
		//t := templ.NewArrayTemplate("${begin|id}", "${end|id}")
		//text := t.Array(context.Background(), template, project.Env, templ.BuildNames(v.Name), ModelArrMapInitF(primaryFields))
		//t = templ.NewArrayTemplate("${begin|no:id}", "${end|no:id}")
		//text = t.Array(context.Background(), text, project.Env, templ.BuildNames(v.Name), ModelArrMapInitF(noPrimaryFields))
		t := templ.NewArrayTemplate("${begin", "${end}", "${case ", "${endcase}")
		text := t.Array(context.Background(), template, project.Env, templ.BuildNames(v.Name), ModelArrMapInitF(v.Fields))
		outputFile = append(outputFile, metadata.File{project.Env["model"] + "/" + v.Name + ".go", text})
	}
	return outputFile
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
func TemplateMapInitF(project metadata.Project) (map[string]string, error) {
	templateMap := make(map[string]string)
	for _, v := range project.Statics {
		content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + v.Name + ".txt")
		if err != nil {
			return nil, err
		}
		templateMap[v.Name] = string(content)
	}
	for _, v := range project.Arrays {
		content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + v.Name + ".txt")
		if err != nil {
			return nil, err
		}
		templateMap[v.Name] = string(content)
	}
	for _, v := range project.Entities {
		content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + v.Name + ".txt")
		if err != nil {
			return nil, err
		}
		templateMap[v.Name] = string(content)
	}
	content, err := ioutil.ReadFile(defaultTemplateFolder + string(os.PathSeparator) + "model.txt")
	if err != nil {
		return nil, err
	}
	templateMap["model"] = string(content)
	return templateMap, nil
}
func InputStructToOutputString(result *string) string {
	//WRITE THE OUT STRUCT
	output.RootPath = defaultRootPath
	output.ProjectName = input.Project.Env["projectName"]
	output.RootPath = strings.TrimSuffix(output.RootPath, "/")
	templateMap, err := TemplateMapInitF(input.Project)
	if err != nil {
		return err.Error()
	}

	var t generator.Generator
	t = &DefaultGenerator{}
	outputFiles := t.Generate(context.Background(), input.Project, templateMap, nil)
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
	var file metadata.File
	for _, k := range content.Files {
		out.OutFile = append(output.OutFile, WriteStruct(&k))
		file.Name = content.Model + "/" + ToLower(k.Name) + ".go"
		file.Content = k.WriteFile.String()
		out.Files = append(out.Files, file)
	}
}
