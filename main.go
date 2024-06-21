package main

import (
	"LPT/data"
	"LPT/resources"
	"LPT/widgets"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var mainApplication fyne.App = app.New()

func main() {
	mainWindow := mainApplication.NewWindow("LPT settings")
	onSettingsFail := func(input string, err error) {
		fmt.Println(input, err)
	}

	particlesPlaceholder := container.NewStack(
		canvas.NewImageFromResource(resources.ResourceNoDataJpg),
		widget.NewLabel("Необходимо указать параметры интерполяции и размер поля"),
	)
	displayPlaceholder := container.NewStack(
		canvas.NewImageFromResource(resources.ResourceNoDataJpg),
		widget.NewLabel("Необходимо указать параметры интерполяции и поле частиц"),
	)
	onSettingsParsed := func(params data.FieldParams, tracing bool, vField, hField []*data.VelocityField) {
		fieldManager, err := data.NewFieldManager(params, tracing, data.NewEmptyField(params.Size), vField, hField)
		if err != nil {
			panic(err)
		}
		particlesPlaceholder.RemoveAll()
		particlesPlaceholder.Add(widgets.NewEditWidget(mainWindow, fieldManager, 720, 720))
		displayPlaceholder.RemoveAll()
		displayPlaceholder.Add(widgets.NewDisplayMenuWidget(fieldManager, 720, 720))
	}
	settingsWindow := widgets.NewSettingsWidget(mainApplication, mainWindow, onSettingsFail, onSettingsParsed) //widgets.NewSettingsMenu(mainApplication, mainWindow, onChanged)

	tabs := container.NewAppTabs(
		container.NewTabItem("Модель", settingsWindow),
		container.NewTabItem("Поле частиц", particlesPlaceholder),
		container.NewTabItem("Результат", displayPlaceholder),
	)
	tabs.SetTabLocation(container.TabLocationLeading)
	mainWindow.SetContent(tabs)
	mainWindow.Show()
	mainApplication.Run()
}
