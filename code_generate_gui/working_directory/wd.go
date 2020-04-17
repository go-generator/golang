package working_directory

import (
	"log"
	"os"
)

func GetWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	return wd
}
