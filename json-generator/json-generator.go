//USAGE INSTRUCTIONS
//Example: "go run main.go input.json rootPath", without filename.json provided, default filename will be "input.json"...
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const (
	defaultFileName    = "input.json"
	defaultRootPath    = ""
	defaultProjectName = "project"
)

var template = map[string]string{
	"service": `package {env}

import (
	"github.com/common-go/search"
	"github.com/common-go/service"
)
type {entity}Service interface {
	service.GenericService
	search.SearchService
}`,
	"impl": `package {env}

import (
	. "evaluation/model"
	. "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"reflect"
)

type {entity}ServiceImpl struct {
	database   *mongo.Database
	collection *mongo.Collection
	*DefaultGenericService
	*DefaultSearchService
}

func New{entity}ServiceImpl(db *mongo.Database, searchResultBuilder SearchResultBuilder) *{entity}ServiceImpl {
	var model {entity}
	modelType := reflect.TypeOf(model)
	collection := "{entityLowerFirstCharacter}"
	mongoService, searchService := NewMongoGenericSearchService(db, modelType, collection, searchResultBuilder, false, "")
	return &{entity}ServiceImpl{db, db.Collection(collection), mongoService, searchService}
}`,
	"controller": `package {env}

import (
	"../handler"
	"../model"
	"../search-model"
	"../service"
	. "github.com/common-go/echo"
	"reflect"
)

type {entity}Controller struct {
	*GenericController
	*SearchController
}


func New{entity}Controller({entityLowerFirstCharacter}Service service.{entity}Service, logService ActivityLogService) *{entity}Controller {
	var {entityLowerFirstCharacter}Model model.{entity}
	modelType := reflect.TypeOf({entityLowerFirstCharacter}Model)
	searchModelType := reflect.TypeOf(search_model.{entity}SM{})
	idNames := {entityLowerFirstCharacter}Service.GetIdNames()
	controlModelHandler := handler.NewControlModelHandler(idNames)
	genericController, searchController:= NewGenericSearchController({entityLowerFirstCharacter}Service, modelType, controlModelHandler, {entityLowerFirstCharacter}Service, searchModelType,nil, logService, true, "")
	return &{entity}Controller{GenericController: genericController, SearchController: searchController}
}`}

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
	for k := range input.Folders {
		for i := range input.Folders[k].RawEnv {
			tmp := strings.LastIndex(input.Folders[k].RawEnv[i], "/")
			input.Folders[k].Env = append(input.Folders[k].Env, input.Folders[k].RawEnv[i][tmp+1:])
			for j := range input.Folders[k].Entity {
				text := template[input.Folders[k].Env[i]]
				text = strings.ReplaceAll(text, "{env}", input.Folders[k].Env[i])
				text = strings.ReplaceAll(text, "{entity}", input.Folders[k].Entity[j])
				text = strings.ReplaceAll(text, "{entityLowerFirstCharacter}", string(strings.ToLower(input.Folders[k].Entity[j])[0])+input.Folders[k].Entity[j][1:])
				filename := FileNameConverter(input.Folders[k].Entity[j], input.Folders[k].RawEnv[i])
				output.Files = append(output.Files, File{input.Folders[k].RawEnv[i] + "/" + filename, text})
			}
		}
	}
	output.RootPath = rootPath
	output.ProjectName = defaultProjectName

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

	////New-style ZIP
	//if output.ProjectName == "" {
	//	output.ProjectName = defaultProjectName
	//}
	//newZipFile, err := os.Create(output.RootPath + output.ProjectName + ".zip")
	//if err != nil {
	//	panic(err)
	//}
	//defer newZipFile.Close()
	//w := zip.NewWriter(newZipFile)
	//for i := range output.Files {
	//	output.Files[i].Name = strings.TrimPrefix(output.Files[i].Name, "/")
	//	f, err := w.Create(output.Files[i].Name)
	//	if err != nil {
	//		panic(err)
	//	}
	//	_, err = f.Write([]byte(output.Files[i].Content))
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	//err = w.Close()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Zip created on disk")

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
