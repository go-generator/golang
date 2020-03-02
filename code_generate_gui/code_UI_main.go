// Package main provides various examples of Fyne API capabilities
package main

import (
	. "./cache"
	"./screens"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
)

func main() {
	a := app.NewWithID("io.fyne.demo")
	a.Settings().SetTheme(theme.LightTheme())
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("Code Generator")
	w.Resize(fyne.Size{
		Width: 640,
	})
	w.SetMaster()
	w.SetContent(screens.WidgetScreen(a, AbsPath))
	w.ShowAndRun()
}
