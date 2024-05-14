package widgets

import (
	"LPT/data"
	"fmt"
	"image"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DisplayWindow struct {
	image        *ImageDisplay
	fieldManager *data.FieldManager
	window       fyne.Window
	currentT     int
	nextPageBtn  *widget.Button
	prevPageBtn  *widget.Button
	playBtn      *widget.Button
	pageLb       binding.String
	maxT         int
	pause        bool
	timerQ       chan struct{}
}

type ImageDisplay struct {
	img *canvas.Image
}

func NewDisplayWindow(app fyne.App, title string, w, h float32, fieldManager *data.FieldManager) *DisplayWindow {
	display := &DisplayWindow{
		window:       app.NewWindow(title),
		fieldManager: fieldManager,
		image:        NewImageDisplay(fieldManager.GetImage(0, 1080, 1080)),
		maxT:         fieldManager.VelocityRecords,
		pageLb:       binding.NewString(),
		pause:        true,
	}
	display.prevPageBtn = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), display.PreviousStep)
	display.nextPageBtn = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), display.NextStep)
	display.playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), display.PlayPause)
	arrows := container.New(layout.NewHBoxLayout(),
		display.prevPageBtn,
		container.New(layout.NewVBoxLayout(), display.playBtn, widget.NewLabelWithData(display.pageLb)),
		display.nextPageBtn,
	)
	canvasContainer := container.NewStack(display.image.GetCanvas())
	canvasContainer.Resize(fyne.NewSize(w, h))
	container := container.New(layout.NewVBoxLayout(), canvasContainer, container.NewCenter(arrows))
	display.prevPageBtn.Disable()
	display.pageLb.Set(fmt.Sprintf("%d/%d", display.currentT+1, display.maxT))
	display.window.SetContent(container)
	return display
}

func (w *DisplayWindow) Show() { w.window.Show() }
func (w *DisplayWindow) Hide() { w.window.Hide() }
func (w *DisplayWindow) PlayPause() {
	if w.pause {
		w.nextPageBtn.Disable()
		w.prevPageBtn.Disable()
		ticker := time.NewTicker(10 * time.Millisecond)
		w.timerQ = make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					if w.currentT >= w.maxT-1 {
						w.currentT = 0
					}
					w.NextStep()
				case <-w.timerQ:
					ticker.Stop()
					return
				}
			}
		}()
		w.pause = false
	} else {
		w.nextPageBtn.Enable()
		w.prevPageBtn.Enable()
		close(w.timerQ)
		w.pause = true
	}
}
func (w *DisplayWindow) PreviousStep() {
	if w.pause && w.currentT < w.maxT-1 {
		w.nextPageBtn.Enable()
	}
	w.currentT--
	img := w.fieldManager.GetImage(w.currentT, 1080, 1080)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT+1, w.maxT))
	if w.currentT == 0 {
		w.prevPageBtn.Disable()
	}
}
func (w *DisplayWindow) NextStep() {
	if w.pause && w.currentT == 0 {
		w.prevPageBtn.Enable()
	}
	w.currentT++
	img := w.fieldManager.GetImage(w.currentT, 1080, 1080)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT+1, w.maxT))
	if w.currentT == w.maxT-1 {
		w.nextPageBtn.Disable()
	}
}

func NewImageDisplay(img image.Image) *ImageDisplay {
	imageDis := &ImageDisplay{}
	imageDis.img = canvas.NewImageFromImage(img)
	imageDis.img.FillMode = canvas.ImageFillOriginal
	return imageDis
}

func (imageDis *ImageDisplay) SetImage(img image.Image) {
	imageDis.img.Image = img
	imageDis.img.Refresh()
}

func (imageDis *ImageDisplay) GetCanvas() *canvas.Image { return imageDis.img }
