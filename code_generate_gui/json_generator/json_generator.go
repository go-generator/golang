package json_generator

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	. "../cache"
	. "../cache/yaml/cache_cipher"
	. "../db_config"
	. "../db_relationship"
	"fyne.io/fyne"
)

const (
	packageFolder = "model"
	outputFile    = "model"
)

func RunWithCommandLine() {
	var dbConfig DatabaseConfig
	err := ReadCacheFile(AbsPath, &dbConfig, "Password")
	if err != nil {
		dialectPtr := flag.String("dialect", "", "input dialect")
		userPtr := flag.String("username", "", "input username")
		passPtr := flag.String("password", "", "input password")
		hostPtr := flag.String("host", "", "input host")
		portPtr := flag.Int("port", 0, "input port")
		dbNamePtr := flag.String("dbName", "", "input database name")
		flag.Parse()
		cacheData := DatabaseConfig{
			Dialect:  *dialectPtr,
			User:     *userPtr,
			Password: *passPtr,
			Host:     *hostPtr,
			Port:     *portPtr,
			Database: *dbNamePtr,
		}
		dbConfig = cacheData
		err = WriteCacheFile(AbsPath, &cacheData, "Password")
		if err != nil {
			log.Fatal(err)
		}
	}
	conn, err := dbConfig.ConnectToSqlServer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = conn.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	rl, _ := DatabaseRelationships(dbConfig, conn)
	err = JsonDescriptionGenerator(packageFolder, outputFile, conn, &dbConfig, rl)
	if err != nil {
		log.Println(err)
	}
}

func RunWithUI(app fyne.App, absPath string) (fyne.Window, error) {
	var dbConfig DatabaseConfig
	var w fyne.Window
	encryptField := "Password"
	if _, err := os.Stat(absPath); os.IsNotExist(err) { // Check if cache file exists, if not then create one
		err = ioutil.WriteFile(absPath, nil, 0666)
		if err != nil {
			ShowWindows(app, "Error", err.Error())
			return w, err
		}
	}
	fs, err := os.Stat(absPath)
	if err != nil {
		ShowWindows(app, "Error", err.Error())
		return w, err
	}
	if fs.Size() != 0 { // Check if cache file is empty
		err = ReadCacheFile(absPath, &dbConfig, encryptField)
		if err != nil {
			ShowWindows(app, "Error", err.Error())
		}
	}
	w = InputUI(&dbConfig, app, absPath, encryptField)
	return w, err
}
