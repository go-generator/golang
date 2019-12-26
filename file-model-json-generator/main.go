package main

import . "./model_json"

func main() {
	source := "./json_input/"
	destination := "./json_output/"
	projectName := ""
	rootPath := ""
	fileOutput := "model_json_output"
	ModelJSONFileGenerator(source, destination, projectName, rootPath, fileOutput)

}
