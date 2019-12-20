//USAGE INSTRUCTIONS
//Example: "go run main.go filename.json", without filename.json provided, default filename will be "input.json"
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
	defaultFileName = "input.json"
)

type Input struct {
	ProjectName string `json:"projectName"`
	RootPath    string `json:"rootPath"`
	Files       []File `json:"files"`
}
type File struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func main() {

	//READ THE JSON FILE

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
	err = json.Unmarshal(byteValue, &input)
	if err != nil {
		panic(err)
	}

	input.RootPath = strings.TrimSuffix(input.RootPath, "/")

	//CREATE FOLDER ON DISK

	err = os.MkdirAll(input.RootPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	for i := range input.Files {
		tmpPath := input.RootPath + "/" + input.Files[i].Name
		tmp := strings.LastIndex(tmpPath, "/")
		tmpPath = tmpPath[:tmp]
		err = os.MkdirAll(tmpPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(input.RootPath + "/" + input.Files[i].Name)
		if err != nil {
			panic(err)
		}
		_, err = f.WriteString(input.Files[i].Content)
		if err != nil {
			panic(err)
		}
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
	fmt.Println(strconv.Itoa(len(input.Files)) + " files created on disk")

	//New-style ZIP

	err = os.MkdirAll(input.RootPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	newZipFile, err := os.Create(input.RootPath + "/" + input.ProjectName + ".zip")
	if err != nil {
		panic(err)
	}
	defer newZipFile.Close()
	w := zip.NewWriter(newZipFile)
	for i := range input.Files {
		input.Files[i].Name = strings.TrimPrefix(input.Files[i].Name, "/")
		f, err := w.Create(input.Files[i].Name)
		if err != nil {
			panic(err)
		}
		_, err = f.Write([]byte(input.Files[i].Content))
		if err != nil {
			panic(err)
		}
	}
	err = w.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Zip created on disk")

}
