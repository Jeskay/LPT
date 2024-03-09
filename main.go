package main

import (
	"LPT/data"
	"bufio"
	"fmt"
	"image"
	"io"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var field *data.Field
var paused bool = true
var mainApplication fyne.App = app.New()

func main() {

	settingsWindow := mainApplication.NewWindow("LPT settings")

	stepSetting := widget.NewEntry()
	particleAmountSettings := widget.NewEntry()
	widthSetting := widget.NewEntry()
	heightSetting := widget.NewEntry()
	velocityButton := widget.NewButton("Выбрать файл v", onVelocityFile(settingsWindow))
	uButton := widget.NewButton("Выбрать файл u", onUFile(settingsWindow))
	wButton := widget.NewButton("Выбрать файл w", onWFile(settingsWindow))
	settingsMenu := widget.NewForm(
		widget.NewFormItem("Шаг", stepSetting),
		widget.NewFormItem("Количество частиц", particleAmountSettings),
		widget.NewFormItem("Ширина", widthSetting),
		widget.NewFormItem("Длина", heightSetting),
		widget.NewFormItem("", velocityButton),
		widget.NewFormItem("", uButton),
		widget.NewFormItem("", wButton),
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
		displayWindowStartup()
		paused = false

	}
	settingsMenu.OnCancel = func() {
		paused = true
	}
	settingsWindow.SetContent(settingsMenu)
	settingsWindow.Resize(fyne.NewSize(1080, 720))

	settingsWindow.Show()

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

func readDataFromFile(r io.Reader) ([][]float64, error) {
	reader := bufio.NewReader(r)
	var numbers [][]float64 = make([][]float64, 256)
	var currentNumber strings.Builder
	row := 0
	counter := 0
	for row < 256 {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		if r == ' ' || r == '\r' || r == '\n' {
			if currentNumber.Len() != 0 {
				num, err := strconv.ParseFloat(currentNumber.String(), 64)
				if err != nil {
					fmt.Println(err)
					return nil, err
				}
				if counter == 256 {
					row++
					counter = 0
				}
				numbers[row] = append(numbers[row], num)
				counter++
				currentNumber.Reset()
			}
			continue
		}
		_, err = currentNumber.WriteRune(r)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return numbers, nil
}

func onVelocityFile(parent fyne.Window) func() {
	return func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err)
			}
			data, err := readDataFromFile(uc)
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
			data, err := readDataFromFile(uc)
			if err != nil {
				return
			}
			fmt.Println(data)
		}, parent)
	}
}

func onWFile(parent fyne.Window) func() {
	return func() {
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				fmt.Println(err)
			}
			data, err := readDataFromFile(uc)
			if err != nil {
				return
			}
			fmt.Println(data)
		}, parent)
	}
}

func displayWindowStartup() {
	w := mainApplication.NewWindow("LPT display")
	raster := canvas.NewRaster(updateImage())
	fmt.Println("field created")

	go func() {
		for range time.Tick(time.Microsecond * 1000) {
			if !paused {
				raster.Refresh()
			}
		}
	}()
	w.SetContent(raster)
	w.Resize(fyne.NewSize(1080, 720))
	w.Show()
}
