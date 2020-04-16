// Package main provides various examples of Fyne API capabilities
package main

import (
	"log"

	. "./cache"
	"./screens"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	//a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("Metadata and Code Generator")
	w.Resize(fyne.Size{
		Width: 640,
	})
	w.SetMaster()
	w.SetContent(screens.WidgetScreen(a, AbsPath))
	w.ShowAndRun()
}
