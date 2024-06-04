package widgets

import (
	"LPT/data"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type DisplayWindow struct {
	fieldManager *data.FieldManager
	window       fyne.Window
}

type DisplayMenuWidget struct {
	widget.BaseWidget
	fieldManager *data.FieldManager
	image        *ImageDisplay
	currentT     int
	nextPageBtn  *widget.Button
	prevPageBtn  *widget.Button
	playBtn      *widget.Button
	pageLb       binding.String
	maxT         int
	pause        bool

	timerQ chan struct{}
}

func NewDisplayMenuWidget(fieldManager *data.FieldManager, width, height float32) *DisplayMenuWidget {
	w := &DisplayMenuWidget{
		pause:        true,
		fieldManager: fieldManager,
		maxT:         fieldManager.VelocityRecords,
		pageLb:       binding.NewString(),
	}
	onDrag := func(i *ImageDisplay, p fyne.Position) {
		const minL float64 = -3
		const maxL float64 = 3
		const rad float64 = 30
		x := (float64(p.X) * (maxL - minL) / float64(i.img.Image.Bounds().Dx())) + minL
		y := (float64(p.Y) * (minL - maxL) / float64(i.img.Image.Bounds().Dy())) - minL
		radius := rad * (maxL - minL) / float64(i.img.Image.Bounds().Dx())
		w.fieldManager.Field.SetColor(x, y, color.RGBA{0, 0, 255, 0xff}, radius)
		img := w.fieldManager.GetCurrentFieldImage(1080, 1080)
		w.image.SetImage(img)
	}
	w.image = NewImageDisplay(fieldManager.GetImageById(0, 1080, 1080), 30, onDrag)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT+1, w.maxT))
	w.prevPageBtn = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), w.PreviousStep)
	w.nextPageBtn = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), w.NextStep)
	w.playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), w.PlayPause)
	w.image.Resize(fyne.NewSize(width, height))
	w.ExtendBaseWidget(w)
	return w
}

func (display *DisplayMenuWidget) CreateRenderer() fyne.WidgetRenderer {
	container := container.New(
		layout.NewVBoxLayout(),
		container.NewStack(display.image),
		container.NewCenter(
			container.New(
				layout.NewHBoxLayout(),
				display.prevPageBtn,
				container.New(
					layout.NewVBoxLayout(),
					display.playBtn,
					widget.NewLabelWithData(display.pageLb)),
				display.nextPageBtn,
			),
		),
	)
	return widget.NewSimpleRenderer(container)
}

func (w *DisplayMenuWidget) PlayPause() {
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
func (w *DisplayMenuWidget) PreviousStep() {
	if w.pause && w.currentT < w.maxT-1 {
		w.nextPageBtn.Enable()
	}
	w.currentT--
	img := w.fieldManager.GetImageById(w.currentT, 1080, 1080)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT+1, w.maxT))
	if w.currentT == 0 {
		w.prevPageBtn.Disable()
	}
}
func (w *DisplayMenuWidget) NextStep() {
	if w.pause && w.currentT == 0 {
		w.prevPageBtn.Enable()
	}
	w.currentT++
	img := w.fieldManager.GetImageById(w.currentT, 1080, 1080)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.currentT+1, w.maxT))
	if w.currentT == w.maxT-1 {
		w.nextPageBtn.Disable()
	}
}

func NewDisplayWindow(app fyne.App, title string, w, h float32, fieldManager *data.FieldManager) *DisplayWindow {
	display := &DisplayWindow{
		window:       app.NewWindow(title),
		fieldManager: fieldManager,
	}
	display.window.SetContent(NewDisplayMenuWidget(fieldManager, w, h))
	return display
}

func (w *DisplayWindow) Show() { w.window.Show() }
func (w *DisplayWindow) Hide() { w.window.Hide() }
