package main

import (
	"LPT/data"
	"LPT/widgets"
	"fmt"

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
	var edit *widgets.EditWidget
	var display *widgets.DisplayMenuWidget
	var editTab *container.TabItem
	var displayTab *container.TabItem
	tabs := container.NewAppTabs()

	onSettingsParsed := func(params data.FieldParams, tracing bool, vField, hField []*data.VelocityField) {
		fieldManager, err := data.NewFieldManager(params, tracing, data.NewEmptyField(params.Size), vField, hField)
		if err != nil {
			panic(err)
		}
		edit = widgets.NewEditWidget(mainWindow, fieldManager, 720, 720)
		editTab = container.NewTabItem("Поле частиц", edit)
		display = widgets.NewDisplayMenuWidget(fieldManager, 720, 720)
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

	mainWindow.SetContent(tabs)
	mainWindow.Show()
	mainApplication.Run()
}
