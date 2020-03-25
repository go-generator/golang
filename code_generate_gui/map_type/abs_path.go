package map_type

import (
	"log"
	"path/filepath"

	"github.com/spf13/viper"
	"golang/code_generate_gui/utils"
	"golang/code_generate_gui/working_directory"
)

var (
	DTypeAbsPath = DataTypeAbsPath()
	sqlTypeMap   map[string]string
)

func DataTypeAbsPath() string {
	filePath := []string{working_directory.GetWorkingDirectory(), "map_type"}
	absPath := filepath.Join(filePath...)
	return absPath
}

func InitTypeMap() {
	var typeMap map[string]string
	viper.SetConfigName("data_type")
	viper.AddConfigPath(DTypeAbsPath)
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Error while reading config file, " + err.Error())
	}
	err := viper.Unmarshal(&sqlTypeMap)
	log.Println(typeMap)
	if err != nil {
		log.Println("Error while unmarshal file, " + err.Error())
	}
}

func RetrieveTypeMap() map[string]string {
	if sqlTypeMap == nil {
		InitTypeMap()
	}
	return utils.CopyMap(sqlTypeMap)
}
