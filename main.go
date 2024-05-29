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
// 	w := mainApplication.NewWindow("DEMO")
// 	img := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{500, 500}})
// 	for i := 0; i < 500; i++ {
// 		for j := 0; j < 500; j++ {
// 			img.Set(i, j, color.RGBA{0, 100, 0, 0xff})
// 		}
// 	}
// 	onDrag := func(imageDisplay *widgets.ImageDisplay, p fyne.Position) {
// 		img.Set(int(p.X), int(p.Y), color.Black)
// 		imageDisplay.SetImage(img)
// 	}
// 	w.SetContent(widgets.NewImageDisplay(img, onDrag))
// 	w.ShowAndRun()
// }

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
