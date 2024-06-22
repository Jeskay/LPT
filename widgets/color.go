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
	colorDis *canvas.Rectangle
	onPick   func(color.RGBA)
	window   fyne.Window
	picker   *dialog.ColorPickerDialog
}

func NewColorPickerWidget(window fyne.Window, onPick func(color.RGBA)) *ColorPickerWidget {
	cPicker := &ColorPickerWidget{
		window:   window,
		onPick:   onPick,
		colorDis: canvas.NewRectangle(color.Black),
	}
	cPicker.colorDis.Resize(fyne.NewSize(cPicker.MinSize().Height, cPicker.MinSize().Height))
	clback := func(clr color.Color) {
		cPicker.colorDis.FillColor = clr
		r, g, b, a := clr.RGBA()
		cPicker.onPick(color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)})
		cPicker.Refresh()
	}
	cPicker.picker = dialog.NewColorPicker("Color", "", clback, cPicker.window)
	cPicker.picker.Advanced = true
	cPicker.ExtendBaseWidget(cPicker)
	return cPicker
}

func (cPicker *ColorPickerWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewGridWithColumns(2,
		container.NewHBox(
			widget.NewLabel("Текущий цвет"),
			widget.NewButton("Изменить", cPicker.onChange),
		),
		cPicker.colorDis,
	)
	return widget.NewSimpleRenderer(c)
}

func (cPicker *ColorPickerWidget) onChange() {
	cPicker.picker.Show()
}
