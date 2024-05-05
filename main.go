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
	settingsWindow := widgets.NewSettingsMenu(mainApplication, mainWindow, onChanged)
	mainWindow.SetContent(settingsWindow.GetForm())
	mainWindow.Show()
	mainApplication.Run()
}

// func main() {
// 	step := 0.09
// 	interSteps := 4
// 	for i := 0; i < 50; i++ {
// 		t := (float64(i*(interSteps+1)) + 1) * step
// 		u, v := utils.GenerateVelocity(t, 3, -3)
// 		utils.WriteDataToFile(u, "uVelocity"+strconv.Itoa(i+1)+".dat")
// 		utils.WriteDataToFile(v, "wVelocity"+strconv.Itoa(i+1)+".dat")
// 	}

// }
