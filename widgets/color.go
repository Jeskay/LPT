package widgets

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type ColorPickerWidget struct {
	widget.BaseWidget
	colorDis *canvas.Circle
	onPick   func(color.RGBA)
	window   fyne.Window
	picker   *dialog.ColorPickerDialog
}

func NewColorPickerWidget(window fyne.Window, onPick func(color.RGBA)) *ColorPickerWidget {
	cPicker := &ColorPickerWidget{
		window:   window,
		onPick:   onPick,
		colorDis: canvas.NewCircle(color.RGBA64{0, 0, 0, 255}),
	}
	clback := func(clr color.Color) {
		cPicker.colorDis.FillColor = clr
		r, g, b, a := clr.RGBA()
		cPicker.onPick(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		cPicker.colorDis.Refresh()
	}
	cPicker.picker = dialog.NewColorPicker("Color", "", clback, cPicker.window)
	cPicker.colorDis.Resize(fyne.NewSize(10, 10))
	cPicker.colorDis.StrokeWidth = 1
	cPicker.colorDis.StrokeColor = color.Black
	cPicker.ExtendBaseWidget(cPicker)
	return cPicker
}

func (cPicker *ColorPickerWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewHBox(
		cPicker.colorDis,
		widget.NewLabel("Текущий цвет"),
		widget.NewButton("Изменить", cPicker.onChange),
	)
	return widget.NewSimpleRenderer(c)
}

func (cPicker *ColorPickerWidget) onChange() {
	cPicker.picker.Show()
}
