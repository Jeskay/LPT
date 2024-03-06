package main

import (
	"LPT/data"
	"fmt"
	"image"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var field *data.Field
var paused bool = true

func main() {
	fmt.Println("starting...")
	a := app.New()
	w := a.NewWindow("Lagrangian particle tracking")

	fmt.Println("creating field...")
	createImage(720, 480, 500, 0.2)
	raster := canvas.NewRaster(updateImage())
	fmt.Println("field created")

	go func() {
		for range time.Tick(time.Microsecond * 1000) {
			if !paused {
				raster.Refresh()
			}
		}
	}()
	stepSetting := widget.NewEntry()
	particleAmountSettings := widget.NewEntry()
	widthSetting := widget.NewEntry()
	heightSetting := widget.NewEntry()
	settingsMenu := widget.NewForm(
		widget.NewFormItem("Шаг", stepSetting),
		widget.NewFormItem("Количество частиц", particleAmountSettings),
		widget.NewFormItem("Ширина", widthSetting),
		widget.NewFormItem("Длина", heightSetting),
	)
	settingsMenu.OnSubmit = func() {
		w, err := strconv.Atoi(widthSetting.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		h, err := strconv.Atoi(heightSetting.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		step, err := strconv.ParseFloat(stepSetting.Text, 64)
		if err != nil {
			fmt.Println(err)
			return
		}
		particles, err := strconv.Atoi(particleAmountSettings.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		paused = true
		createImage(w, h, particles, step)
		paused = false
	}
	settingsMenu.OnCancel = func() {
		paused = true
	}
	grid := container.New(layout.NewGridLayout(2), raster, settingsMenu)
	w.SetContent(grid)
	w.Resize(fyne.NewSize(800, 800))
	w.ShowAndRun()
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
