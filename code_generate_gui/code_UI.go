// Package main provides various examples of Fyne API capabilities
package main

import (
	"net/url"

	. "./cache"
	"./data"
	"./screens"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

const (
	preferenceCurrentTab = "currentTab"
)

func welcomeScreen(a fyne.App) fyne.CanvasObject {
	logo := canvas.NewImageFromResource(data.FyneScene)
	logo.SetMinSize(fyne.NewSize(228, 167))

	link, err := url.Parse("https://fyne.io/")
	if err != nil {
		fyne.LogError("Could not parse URL", err)
	}

	return widget.NewVBox(
		widget.NewLabelWithStyle("Welcome to the Fyne toolkit demo app", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		layout.NewSpacer(),
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewHyperlinkWithStyle("fyne.io", link, fyne.TextAlignCenter, fyne.TextStyle{}),
		layout.NewSpacer(),

		widget.NewGroup("Theme",
			fyne.NewContainerWithLayout(layout.NewGridLayout(2),
				widget.NewButton("Dark", func() {
					a.Settings().SetTheme(theme.DarkTheme())
				}),
				widget.NewButton("Light", func() {
					a.Settings().SetTheme(theme.LightTheme())
				}),
			),
		),
	)
}

func main() {
	a := app.NewWithID("io.fyne.demo")
	a.Settings().SetTheme(theme.LightTheme())
	a.SetIcon(theme.FyneLogo())
	w := a.NewWindow("Code Generator")
	w.Resize(fyne.Size{
		Width: 640,
	})
	//w.SetMainMenu(fyne.NewMainMenu(fyne.NewMenu("File",
	//	fyne.NewMenuItem("New", func() { fmt.Println("Menu New") }),
	//	// a quit item will be appended to our first menu
	//), fyne.NewMenu("Edit",
	//	fyne.NewMenuItem("Cut", func() { fmt.Println("Menu Cut") }),
	//	fyne.NewMenuItem("Copy", func() { fmt.Println("Menu Copy") }),
	//	fyne.NewMenuItem("Paste", func() { fmt.Println("Menu Paste") }),
	//)))
	w.SetMaster()
	//tabs := widget.NewTabContainer(
	//
	//	widget.NewTabItemWithIcon("Widgets", theme.ContentCopyIcon(), screens.WidgetScreen()),
	//)
	//
	//tabs.SetTabLocation(widget.TabLocationLeading)
	//tabs.SelectTabIndex(a.Preferences().Int(preferenceCurrentTab))
	w.SetContent(screens.WidgetScreen(a, AbsPath))
	w.ShowAndRun()
	//a.Preferences().SetInt(preferenceCurrentTab, tabs.CurrentTabIndex())
}
