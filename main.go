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

// func main() {
// 	u, v := utils.GenerateVelocity(3, -3)
// 	for i := 0; i < 50; i++ {
// 		utils.WriteDataToFile(u, "uVelocity"+strconv.Itoa(i)+".dat")
// 		utils.WriteDataToFile(v, "wVelocity"+strconv.Itoa(i)+".dat")
// 	}

// }
