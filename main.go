package main

import (
	"LPT/data"
	"LPT/widgets"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var mainApplication fyne.App = app.New()

func main() {
	mainWindow := mainApplication.NewWindow("LPT settings")
	var displayWindow *widgets.DisplayWindow
	onChanged := func(manager *data.FieldManager) {
		if displayWindow != nil {
			displayWindow.Hide()
		}
		displayWindow = widgets.NewDisplayWindow(mainApplication, "Particle Tracking", 1080, 720, manager)
		displayWindow.Show()
	}
	settingsWindow := widgets.NewSettingsMenu(mainWindow, onChanged)
	mainWindow.SetContent(settingsWindow.GetForm())
	mainWindow.Show()
	mainApplication.Run()
}
