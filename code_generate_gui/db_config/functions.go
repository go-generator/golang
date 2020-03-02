package db_config

import (
	"log"
	"regexp"
	"strings"
	"unicode"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"github.com/jinzhu/gorm"
)

func AddJSONTag(name string) string {
	return "`json:\"" + ToLower(name) + "\""
}

func AddBSONTag(name string) string {
	return " bson:\"" + ToLower(name) + "\""
}

func AddGormTag(name string, primaryTag bool) string {
	if name == "" {
		return "`\n"
	}
	if primaryTag {
		return " gorm:\"column:" + name + ":primary_key\"`\n"
	}
	return " gorm:\"column:" + name + "\"`\n"
}

func AddStructFieldName(name string) string {
	return strings.Title(name)
}

func ToLower(s string) string {
	if len(s) < 0 {
		return ""
	}
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func ListAllTableNames(conn *gorm.DB, databaseName string) []string {
	var res []string
	conn.Table("information_schema.tables").Select("*").Where("table_schema='"+databaseName+"'").Pluck("table_name", &res)
	return res
}

func StandardizeName(s string) string {
	var field strings.Builder
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	tokens := strings.Split(s, "_")
	for _, t := range tokens {
		alphanumericString := reg.ReplaceAllString(t, "")
		field.WriteString(strings.Title(alphanumericString))
	}
	return field.String()
}

func ShowWindows(app fyne.App, title, message string) {
	wa := app.NewWindow(title)
	wa.Resize(fyne.Size{
		Width: 320,
	})
	wa.SetContent(widget.NewVBox(
		widget.NewLabel(message),
	))
	wa.CenterOnScreen()
	wa.Show()
}
