package screens

import (
	"fmt"
	"time"

	"../code_generate_core"
	"../json_generator"
	"github.com/sqweek/dialog"

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
			//largeText2.SetText(result)
			//ShowWindows(app, "Output", largeText2.Text)
			cursorRow.SetText("OK")
		} else {
			cursorRow.SetText(err)
		}
	})
	openFileButton := widget.NewButton("Generate Code From File...", func() {
		filename, err := dialog.File().Filter("JSON/Text file", "json", "txt").Load()
		if err != nil {
			cursorRow.SetText(err.Error())
		} else {
			result := ""
			err := code_generate_core.GenerateFromFile(templateDir.Text, projectName.Text, filename, &result)
			if err == "" {
				largeText2.SetText(result)
				cursorRow.SetText("OK")
			} else {
				cursorRow.SetText(err)
				return
			}
		}
		directory, err := dialog.Directory().Title("Save Generated Files In...").Browse()
		err1 := code_generate_core.OutputStructToFiles(directory)
		if err1 != "" {
			cursorRow.SetText(err1)
		} else {
			cursorRow.SetText("Files Created On Disk")
		}
	})

	saveButton := widget.NewButton("Save Files (to main.go folder or input.json folder)", func() {
		err := code_generate_core.OutputStructToFiles("")
		if err != "" {
			cursorRow.SetText(err)
		} else {
			cursorRow.SetText("Files Created On Disk")
		}
	})
	//saveAsButton := widget.NewButton("Save Files As...", func() {
	//	directory, err := dialog.Directory().Title("Save Files As...").Browse()
	//	if err != nil {
	//		cursorRow.SetText(err.Error())
	//	} else {
	//		err := code_generate_core.OutputStructToFiles(directory)
	//		if err != "" {
	//			cursorRow.SetText(err)
	//		} else {
	//			cursorRow.SetText("Files Created On Disk")
	//		}
	//	}
	//})
	zipButton := widget.NewButton("Save Zip (to main.go folder or input.json folder)", func() {
		err := code_generate_core.OutputStructToZip("")
		if err != "" {
			cursorRow.SetText(err)
		} else {
			cursorRow.SetText("Zip Created On Disk")
		}
	})
	zipAsButton := widget.NewButton("Save Zip As...", func() {
		directory, err := dialog.File().Filter("ZIP files", "zip").Title("Export to ZIP").Save()
		//directory, err := dialog.Directory().Title("Save Zip As...").Browse()
		if err != nil {
			cursorRow.SetText(err.Error())
		} else {
			err := code_generate_core.OutputStructToZip(directory)
			if err != "" {
				cursorRow.SetText(err)
			} else {
				cursorRow.SetText("Zip Created On Disk")
			}
		}
	})
	modelJsonGenerator := widget.NewButton("Json Description Generator", func() {
		wi, err := json_generator.RunWithUI(app, cachePath)
		if err == nil {
			wi.Show()
		}
	})
	list := widget.NewVBox()
	list2 := widget.NewVBox()
	list.Append(templateBar)
	list.Append(okButton)
	list.Append(saveButton)
	list.Append(zipButton)
	list.Append(openFileButton)
	//list.Append(saveAsButton)
	list.Append(zipAsButton)
	list2.Append(projectName)
	list2.Append(largeText)
	//list2.Append(largeText2)
	list.Append(modelJsonGenerator)
	statusBar := widget.NewHBox(layout.NewSpacer(),
		widget.NewLabel("Status:"), cursorRow,
	)
	list.Append(statusBar)

	//	OnCancel: func() {
	//		fmt.Println("Cancelled")
	//	},
	//	OnSubmit: func() {
	//		fmt.Println("Input:", largeText.Text)
	//		fmt.Println("Output:", largeText2.Text)
	//		output:= CodeGenerate(largeText.Text)
	//		largeText2.SetText(output)
	//	},
	//}
	//
	//form.Append("Input", largeText)
	//form.Append("Output", largeText2)
	//scroll := widget.NewScrollContainer(list)
	//scroll.Resize(fyne.NewSize(200, 100))

	scroll2 := widget.NewScrollContainer(list2)
	//scroll2.Resize(fyne.NewSize(200, 100))

	//return form
	return fyne.NewContainerWithLayout(layout.NewGridLayout(1), scroll2, list)
	//return fyne.NewContainerWithLayout(scroll2,list,statusBar)
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
	//toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.MailComposeIcon(), func() { fmt.Println("New") }),
	//	widget.NewToolbarSeparator(),
	//	widget.NewToolbarSpacer(),
	//	widget.NewToolbarAction(theme.ContentCutIcon(), func() { fmt.Println("Cut") }),
	//	widget.NewToolbarAction(theme.ContentCopyIcon(), func() { fmt.Println("Copy") }),
	//	widget.NewToolbarAction(theme.ContentPasteIcon(), func() { fmt.Println("Paste") }),
	//)

	//return fyne.NewContainerWithLayout(layout.NewBorderLayout(nil, nil, nil, nil),
	//	widget.NewTabContainer(
	//		//widget.NewTabItem("Buttons", makeButtonTab()),
	//		//widget.NewTabItem("Input", makeInputTab()),
	//		//widget.NewTabItem("Progress", makeProgressTab()),
	//		widget.NewTabItem("Code Generator", makeFormTab()),
	//		//widget.NewTabItem("Scroll", makeScrollTab()),
	//		//widget.NewTabItem("Full Scroll", makeScrollBothTab()),
	//	),
	//)
	return makeFormTab(app, cachePath)
}
