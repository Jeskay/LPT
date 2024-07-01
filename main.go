package main

import (
	"LPT/data"
	"LPT/widgets"
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
)

var mainApplication fyne.App = app.New()

func main() {
	mainWindow := mainApplication.NewWindow("LPT settings")
	onSettingsFail := func(input string, err error) {
		fmt.Println(input, err)
	}
	var edit *widgets.EditMenu
	var display *widgets.DisplayMenu
	var editTab *container.TabItem
	var displayTab *container.TabItem
	tabs := container.NewAppTabs()

	onSettingsParsed := func(params data.FieldParams, tracing bool, vField, hField []*data.VelocityField) {
		fieldManager, err := data.NewFieldManager(params, tracing, data.NewEmptyField(params.Size), vField, hField)
		if err != nil {
			panic(err)
		}
		if editTab != nil {
			tabs.Remove(editTab)
		}
		if displayTab != nil {
			tabs.Remove(displayTab)
		}
		edit = widgets.NewEditWidget(mainWindow, fieldManager, 720, 720)
		editTab = container.NewTabItem("Поле частиц", edit)
		display = widgets.NewDisplayMenuWidget(func(index int) image.Image {
			return fieldManager.GetImageById(index, 720, 720)
		}, fieldManager.VelocityRecords)
		displayTab = container.NewTabItem("Результат", display)
		tabs.Append(editTab)
		tabs.Append(displayTab)
		tabs.OnSelected = func(ti *container.TabItem) {
			if ti == displayTab {
				fmt.Println("updating")
				fieldManager.ClearHistory()
				display.Update()
			}
		}
	}
	settingsWindow := widgets.NewSettingsWidget(mainApplication, mainWindow, onSettingsFail, onSettingsParsed)

	tabs.Append(container.NewTabItem("Модель", settingsWindow))
	tabs.SetTabLocation(container.TabLocationLeading)
	mainWindow.Resize(fyne.NewSize(1280, 720))
	mainWindow.SetContent(tabs)
	mainWindow.Show()
	mainApplication.Run()
}
