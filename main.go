package main

import (
	"LPT/data"
	"LPT/utils"
	"LPT/widgets"
	"fmt"
	"image"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
)

var field *data.Field
var mainApplication fyne.App = app.New()

func main() {
	mainWindow := mainApplication.NewWindow("LPT settings")
	var displayWindow *widgets.DisplayWindow
	onChanged := func(field *data.Field) {
		if displayWindow != nil {
			displayWindow.Hide()
		}
		displayWindow = widgets.NewDisplayWindow(mainApplication, "Particle Tracking", 1080, 720, field)
		displayWindow.Show()
	}
	settingsWindow := widgets.NewSettingsMenu(mainWindow, onChanged)
	mainWindow.SetContent(settingsWindow.GetForm())
	mainWindow.Show()
	mainApplication.Run()
}

func createImage(w, h int, particleAmount int, step float64) image.Image {
	field = data.NewField(particleAmount, data.Size{
		MinAxisX: -3,
		MaxAxisX: 3,
		MinAxisY: -3,
		MaxAxisY: 3,
	}, step)
	field.UpdatePosition()
	field.Print()
	img := field.Image(w, h)
	return img
}

func updateImage() func(int, int) image.Image {
	return func(w, h int) image.Image {
		field.UpdatePosition()
		return field.Image(w, h)
	}
}

func onVelocityFile(parent fyne.Window) func() {
	return func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err)
			}
			data, err := utils.ReadDataFromFile(uc)
			if err != nil {
				return
			}
			fmt.Println(data)
		}, parent)
	}
}

func onUFile(parent fyne.Window) func() {
	return func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err)
			}
			data, err := utils.ReadDataFromFile(uc)
			if err != nil {
				return
			}
			fmt.Println(data)
		}, parent)
	}
}

func onWFile(parent fyne.Window) func() {
	return func() {
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			readers, err := lu.List()
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, uc := range readers {
				fmt.Println(uc.Path(), "\n", uc.Query())
				f, err := os.Open(uc.Path())
				if err != nil {
					fmt.Println(err)
					continue
				}
				data, err := utils.ReadDataFromFile(f)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fmt.Println(data)
			}
		}, parent)
	}
}
