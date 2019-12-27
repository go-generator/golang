package main

import (
	"flag"

	. "./model_json"
)

func main() {
	source := flag.String("source", "", "input source directory")
	destination := flag.String("destination", "", "input destination directory")
	projectName := flag.String("projectName", "", "input project name")
	rootPath := flag.String("rootPath", "", "input root path")
	output := flag.String("output", "", "input file output name")
	if *output == "" {
		*output = "json_converted"
	}
	if *rootPath == "" {
		*rootPath = ""
	}
	if *projectName == "" {
		*projectName = ""
	}
	ModelJSONFileGenerator(*source, *destination, *projectName, *rootPath, *output)
}
