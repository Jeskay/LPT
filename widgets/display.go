package widgets

import (
	"LPT/data"
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type DisplayWindow struct {
	fieldManager *data.FieldManager
	window       fyne.Window
}

type DisplayMenuWidget struct {
	widget.BaseWidget
	fieldManager    *data.FieldManager
	image           *ImageDisplay
	spdSlider       *IntSlider
	pageProgressBar *widget.ProgressBar
	currentT        int
	nextPageBtn     *widget.Button
	prevPageBtn     *widget.Button
	playBtn         *PauseWidget
	pageLb          binding.String
	fps             int
	fpsStr          binding.String
	maxT            int
	pause           bool

	timerQ chan struct{}
}

func NewDisplayMenuWidget(fieldManager *data.FieldManager, width, height float32) *DisplayMenuWidget {
	w := &DisplayMenuWidget{
		pause:        true,
		fieldManager: fieldManager,
		pageLb:       binding.NewString(),
		fpsStr:       binding.NewString(),
		fps:          1,
	}
	w.fpsStr.Set(strconv.Itoa(w.fps) + "x")
	w.maxT = fieldManager.VelocityRecords
	w.image = NewImageDisplay(fieldManager.GetImageById(0, 1080, 1080), 30, nil)
	w.pageLb.Set(fmt.Sprintf("Кадр %d/%d", w.currentT+1, w.maxT))
	w.prevPageBtn = widget.NewButton("<<Предыдущий", w.PreviousStep)
	w.nextPageBtn = widget.NewButton("Следующий>>", w.NextStep)
	w.pageProgressBar = widget.NewProgressBar()
	w.pageProgressBar.Min = 0
	w.pageProgressBar.Max = float64(w.maxT - 1)
	w.pageProgressBar.TextFormatter = func() string { return "" }
	w.spdSlider = NewIntSlider(w.onSpdChange, w.fps, 120)
	w.playBtn = NewPauseWidget(w.PlayPause)
	w.image.Resize(fyne.NewSize(width, height))
	w.playBtn.Resize(fyne.NewSize(20, 20))
	w.playBtn.Refresh()

	w.ExtendBaseWidget(w)
	return w
}

func (display *DisplayMenuWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.New(
		layout.NewVBoxLayout(),
		container.NewCenter(container.NewStack(display.image, display.playBtn)),
		container.NewBorder(nil, nil,
			widget.NewLabel("Кадров в секунду"),
			widget.NewLabelWithData(display.fpsStr),
			display.spdSlider,
		),
		container.NewPadded(
			container.NewVBox(
				widget.NewLabelWithData(display.pageLb),
				display.pageProgressBar,
			),
		),
		container.NewCenter(
			container.New(
				layout.NewHBoxLayout(),
				display.prevPageBtn,
				display.nextPageBtn,
			),
		),
	)
	return widget.NewSimpleRenderer(c)
}

func (w *DisplayMenuWidget) PlayPause() {
	if w.pause {
		w.nextPageBtn.Disable()
		w.prevPageBtn.Disable()
		ticker := time.NewTicker(time.Duration(1000/w.fps) * time.Millisecond)
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

func (w *DisplayMenuWidget) onSpdChange(value int) {
	w.fps = value
	w.fpsStr.Set(strconv.Itoa(w.fps) + "x")
}

func (w *DisplayMenuWidget) PreviousStep() {
	if w.pause && w.currentT < w.maxT-1 {
		w.nextPageBtn.Enable()
	}
	w.currentT--
	w.pageProgressBar.SetValue(float64(w.currentT))
	img := w.fieldManager.GetImageById(w.currentT, 1080, 1080)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("Кадр %d/%d", w.currentT+1, w.maxT))
	if w.currentT == 0 {
		w.prevPageBtn.Disable()
	}
}
func (w *DisplayMenuWidget) NextStep() {
	if w.pause && w.currentT == 0 {
		w.prevPageBtn.Enable()
	}
	w.currentT++
	w.pageProgressBar.SetValue(float64(w.currentT))
	img := w.fieldManager.GetImageById(w.currentT, 1080, 1080)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("Кадр %d/%d", w.currentT+1, w.maxT))
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
