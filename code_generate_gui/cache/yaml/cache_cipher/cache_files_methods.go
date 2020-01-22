package cache_cipher

import (
	"io/ioutil"
	"log"
	math "math/rand"
	"reflect"
	"strings"
	"time"

	. "./cipher"
	"gopkg.in/yaml.v2"
)

const (
	key = "xbmcZMpQoGiRXlTSbHSuYPXynluuNyYh"
)

func ErrorHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Generate32ByteString() string {
	var output string
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	math.Seed(time.Now().UnixNano())
	for i := 0; i < 32; i++ {
		index := math.Intn(1000)
		if -1 < index && index < 52 {
			output += string(letters[index])
		} else {
			output += string(letters[index%52])
		}
	}
	return output
}

func GetFieldStringValue(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func SetFieldStringValue(inter interface{}, field string, value string) {
	ps := reflect.ValueOf(inter)
	// struct
	e := ps.Elem()
	if e.Kind() == reflect.Struct {
		// exported field
		f := e.FieldByName(field)
		if f.IsValid() {
			// A Value can be changed only if it is
			// addressable and was not obtained by
			// the use of unexported struct fields.
			if f.CanSet() {
				if f.Kind() == reflect.String {
					f.SetString(value)
				}
			}
		}
	}
}

func ReadCacheFile(filePath string, outer interface{}, field string) error {
	var s strings.Builder
	in, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(in, outer)
	if err != nil {
		return err
	}
	plain, err := Decrypt([]byte(key), []byte(GetFieldStringValue(outer, field)))
	if err != nil {
		return err
	}
	s.Write(plain)
	SetFieldStringValue(outer, field, s.String())
	s.Reset()
	return err
}

func WriteCacheFile(filePath string, inter interface{}, field string) error {
	var s strings.Builder
	ciphered, err := Encrypt([]byte(key), []byte(GetFieldStringValue(inter, field)))
	s.Write(ciphered)
	if err != nil {
		return err
	}
	SetFieldStringValue(inter, field, s.String())
	s.Reset()
	data, err1 := yaml.Marshal(inter)
	if err1 != nil {
		return err1
	}
	err = ioutil.WriteFile(filePath, data, 0666)
	return err
}

//func main() {
//	filePath := "./cache_file/cache.yaml"
//	encryptField := "Password"
//	dbConfig := DatabaseConfig{
//		Dialect:  "mysql",
//		Host:     "localhost",
//		Port:     3306,
//		Database: "classicmodels",
//		User:     "test",
//		Password: "Doraemon1096~",
//	}
//	err := WriteCacheFile(filePath, &dbConfig, encryptField)
//	ErrorHandle(err)
//	err = ReadCacheFile(filePath, &dbConfig, encryptField)
//	ErrorHandle(err)
//	log.Println(dbConfig)
//}
