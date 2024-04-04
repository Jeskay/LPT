package widgets

import (
	"LPT/data"
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DisplayWindow struct {
	image       *ImageDisplay
	window      fyne.Window
	currentT    int
	nextPageBtn *widget.Button
	prevPageBtn *widget.Button
	pageLb      binding.String
	maxT        int
}

type ImageDisplay struct {
	img   *canvas.Image
	field *data.Field
}

func NewDisplayWindow(app fyne.App, title string, w, h float32, field *data.Field) *DisplayWindow {
	display := &DisplayWindow{
		window: app.NewWindow(title),
		image:  NewImageDisplay(field),
		maxT:   10,
		pageLb: binding.NewString(),
	}
	display.prevPageBtn = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), display.PreviousStep)
	display.nextPageBtn = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), display.NextStep)
	arrows := container.New(layout.NewHBoxLayout(),
		display.prevPageBtn,
		widget.NewLabelWithData(display.pageLb),
		display.nextPageBtn,
	)
	canvasContainer := container.NewStack(display.image.GetCanvas())
	canvasContainer.Resize(fyne.NewSize(w, h))
	container := container.New(layout.NewVBoxLayout(), canvasContainer, container.NewCenter(arrows))
	display.NextStep()
	display.prevPageBtn.Disable()
	display.window.SetContent(container)
	return display
}

func (w *DisplayWindow) Show() { w.window.Show() }
func (w *DisplayWindow) Hide() { w.window.Hide() }
func (w *DisplayWindow) PreviousStep() {
	if w.currentT == w.maxT {
		w.nextPageBtn.Enable()
	}
	w.currentT--
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT, w.maxT))
	if w.currentT == 1 {
		w.prevPageBtn.Disable()
	}
}
func (w *DisplayWindow) NextStep() {
	if w.currentT == 1 {
		w.prevPageBtn.Enable()
	}
	w.currentT++
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT, w.maxT))
	if w.currentT == w.maxT {
		w.nextPageBtn.Disable()
	}
}

func NewImageDisplay(field *data.Field) *ImageDisplay {
	imageDis := &ImageDisplay{field: field}
	imageDis.img = canvas.NewImageFromImage(imageDis.updateImage(1080, 720))
	imageDis.img.FillMode = canvas.ImageFillOriginal
	return imageDis
}

func (imageDis *ImageDisplay) Refresh() {
	imageDis.img.Image = imageDis.updateImage(1080, 720)
	imageDis.img.Refresh()
}
func (imageDis *ImageDisplay) GetCanvas() *canvas.Image { return imageDis.img }

func (imageDis *ImageDisplay) updateImage(w, h int) image.Image {
	imageDis.field.UpdatePosition()
	return imageDis.field.Image(w, h)
}
