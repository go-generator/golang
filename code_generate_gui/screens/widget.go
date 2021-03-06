package screens

import (
	"fmt"
	"github.com/go-generator/metadata"
	"time"

	"../json_generator"
	"github.com/sqweek/dialog"
	"golang/code_generate_gui/code_generate_core"
	"golang/code_generate_gui/db_config"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func makeButtonTab() fyne.Widget {
	disabled := widget.NewButton("Disabled", func() {})
	disabled.Disable()

	return widget.NewVBox(
		widget.NewLabel("Text label"),
		widget.NewButton("Text button", func() { fmt.Println("tapped text button") }),
		widget.NewButtonWithIcon("With icon", theme.ConfirmIcon(), func() { fmt.Println("tapped icon button") }),
		disabled,
	)
}

func makeInputTab() fyne.Widget {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Entry")
	entryReadOnly := widget.NewEntry()
	entryReadOnly.SetText("Entry (disabled)")
	entryReadOnly.Disable()

	disabledCheck := widget.NewCheck("Disabled check", func(bool) {})
	disabledCheck.Disable()
	radio := widget.NewRadio([]string{"Radio Item 1", "Radio Item 2"}, func(s string) { fmt.Println("selected", s) })
	radio.Horizontal = true
	disabledRadio := widget.NewRadio([]string{"Disabled radio"}, func(string) {})
	disabledRadio.Disable()

	return widget.NewVBox(
		entry,
		entryReadOnly,
		widget.NewSelect([]string{"Option 1", "Option 2", "Option 3"}, func(s string) { fmt.Println("selected", s) }),
		widget.NewCheck("Check", func(on bool) { fmt.Println("checked", on) }),
		disabledCheck,
		radio,
		disabledRadio,
		widget.NewSlider(0, 100),
	)
}

func makeProgressTab() fyne.Widget {
	progress := widget.NewProgressBar()
	infProgress := widget.NewProgressBarInfinite()

	go func() {
		num := 0.0
		for num < 1.0 {
			time.Sleep(100 * time.Millisecond)
			progress.SetValue(num)
			num += 0.01
		}

		progress.SetValue(1)
	}()

	return widget.NewVBox(
		widget.NewLabel("Percent"), progress,
		widget.NewLabel("Infinite"), infProgress)
}

func makeFormTab(app fyne.App, cachePath string) fyne.CanvasObject {
	projectName := widget.NewEntry()
	projectName.SetPlaceHolder("Project Name")
	templateDir := widget.NewEntry()
	templateDir.SetPlaceHolder("Template Folder Directory")
	templateButton := widget.NewButton("Browse...", func() {
		directory, err := dialog.Directory().Title("Browse...").Browse()
		if err != nil {
			templateDir.SetText(err.Error())
			return
		}
		templateDir.SetText(directory)
	})
	templateBar := widget.NewHBox(
		templateDir, templateButton,
	)
	largeText := widget.NewMultiLineEntry()
	largeText.SetPlaceHolder("Input")
	largeText2 := widget.NewMultiLineEntry()
	largeText2.SetPlaceHolder("Output")
	cursorRow := widget.NewLabel("")
	okButton := widget.NewButton("Code Generate", func() {
		result := ""
		err := code_generate_core.GenerateFromString(templateDir.Text, projectName.Text, largeText2.Text, &result)
		if err == "" {
			cursorRow.SetText("OK")
		} else {
			cursorRow.SetText(err)
		}
	})
	var output metadata.Output
	openFileButton := widget.NewButton("Generate Go Template From Metadata Json...", func() {
		filename, err := dialog.File().Filter("json file", "json").Load()
		if err != nil {
			cursorRow.SetText(err.Error())
			return
		} else {
			result := ""
			err := ""
			output, err = code_generate_core.GenerateFromFile(templateDir.Text, projectName.Text, filename, &result, "java")
			if err == "" {
				largeText2.SetText(result)
				cursorRow.SetText("Creating Files...")
			} else {
				cursorRow.SetText(err)
				return
			}
		}
		directory, err := dialog.Directory().Title("Save Generated Files In...").Browse()
		err1 := code_generate_core.OutputStructToFiles(directory, output)
		if err1 != "" {
			cursorRow.SetText(err1)
		} else {
			cursorRow.SetText("Files Created On Disk")
		}
	})

	saveButton := widget.NewButton("Save Files (to main.go folder or input.json folder)", func() {
		err := code_generate_core.OutputStructToFiles("", output)
		if err != "" {
			cursorRow.SetText(err)
		} else {
			cursorRow.SetText("Files Created On Disk")
		}
	})
	//zipButton := widget.NewButton("Save Zip (to main.go folder or input.json folder)", func() {
	//	err := code_generate_core.OutputStructToZip()
	//	if err != "" {
	//		cursorRow.SetText(err)
	//	} else {
	//		cursorRow.SetText("Zip Created On Disk")
	//	}
	//})
	//zipAsButton := widget.NewButton("Generate And Zip From Metadata Json As...", func() {
	//	cursorRow.SetText("Creating Zip File...")
	//	err := code_generate_core.OutputStructToZip()
	//	if err != "" {
	//		cursorRow.SetText(err)
	//	} else {
	//		cursorRow.SetText("Zip Created On Disk")
	//	}
	//})
	modelJsonGenerator := widget.NewButton("Generate Database Metadata Json Description", func() {
		wi, err := json_generator.RunWithUI(app, cachePath)
		if err == nil {
			wi.Show()
		} else {
			cursorRow.SetText(err.Error())
		}
	})
	//TODO: Add Java generator
	javaFilesGenerator := widget.NewButton("Generate Java Files from Json...", func() {
		err := json_generator.JavaFilesGenerator(app, cachePath)
		if err == nil {
			db_config.ShowWindows(app, "Success", "Generated Java Files Successfully")
		} else {
			cursorRow.SetText(err.Error())
		}
	})
	list := widget.NewVBox()
	list2 := widget.NewVBox()
	list.Append(templateBar)
	list.Append(okButton)
	list.Append(saveButton)
	list.Append(openFileButton)
	//list.Append(zipAsButton)
	list2.Append(projectName)
	list2.Append(largeText)
	list.Append(modelJsonGenerator)
	list.Append(javaFilesGenerator)
	statusBar := widget.NewHBox(layout.NewSpacer(),
		widget.NewLabel("Status:"), cursorRow,
	)
	list.Append(statusBar)
	scroll2 := widget.NewScrollContainer(list2)
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1), scroll2, list)
}

func makeScrollTab() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(320, 320))
	list := widget.NewHBox()
	list2 := widget.NewVBox()
	for i := 1; i <= 20; i++ {
		index := i // capture
		list.Append(widget.NewButton(fmt.Sprintf("Button %d", index), func() {
			fmt.Println("Tapped", index)
		}))
		list2.Append(widget.NewButton(fmt.Sprintf("Button %d", index), func() {
			fmt.Println("Tapped", index)
		}))
	}
	scroll := widget.NewScrollContainer(list)
	scroll.Resize(fyne.NewSize(200, 300))
	scroll2 := widget.NewScrollContainer(list2)
	scroll2.Resize(fyne.NewSize(200, 100))
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1), scroll, scroll2)
}

func makeScrollBothTab() fyne.CanvasObject {
	logo := canvas.NewImageFromResource(theme.FyneLogo())
	logo.SetMinSize(fyne.NewSize(800, 800))

	scroll := widget.NewScrollContainer(logo)
	scroll.Resize(fyne.NewSize(400, 400))

	return scroll
}

// WidgetScreen shows a panel containing widget demos
func WidgetScreen(app fyne.App, cachePath string) fyne.CanvasObject {
	return makeFormTab(app, cachePath)
}
