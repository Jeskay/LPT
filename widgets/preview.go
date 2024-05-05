package widgets

import (
	"LPT/data"
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type PreviewWindow struct {
	image   *ImageDisplay
	fields  []*data.VelocityField
	pause   bool
	current int
	window  fyne.Window
	nextBtn *widget.Button
	prevBtn *widget.Button
	playBtn *widget.Button
	pageLb  binding.String
	timerQ  chan struct{}
}

func NewPreviewWindow(app fyne.App, title string, w, h float32, fields []*data.VelocityField) *PreviewWindow {
	preview := &PreviewWindow{
		window: app.NewWindow(title),
		fields: fields,
		image:  NewImageDisplay(data.NewGradientFromField(fields[0], 720, 720)),
		pause:  true,
		pageLb: binding.NewString(),
	}
	preview.prevBtn = widget.NewButtonWithIcon("", theme.MediaFastRewindIcon(), preview.Previous)
	preview.nextBtn = widget.NewButtonWithIcon("", theme.MediaFastForwardIcon(), preview.Next)
	preview.playBtn = widget.NewButtonWithIcon("", theme.MediaPlayIcon(), preview.PlayPause)
	arrows := container.New(layout.NewHBoxLayout(),
		preview.prevBtn,
		container.New(layout.NewVBoxLayout(), preview.playBtn, widget.NewLabelWithData(preview.pageLb)),
		preview.nextBtn,
	)
	canvasContainer := container.NewStack(preview.image.GetCanvas())
	canvasContainer.Resize(fyne.NewSize(w, h))
	container := container.New(layout.NewVBoxLayout(), canvasContainer, container.NewCenter(arrows))
	preview.prevBtn.Disable()
	preview.pageLb.Set(fmt.Sprintf("%d/%d", preview.current+1, len(preview.fields)))
	preview.window.SetContent(container)
	return preview
}

func (w *PreviewWindow) Next() {
	if w.current == len(w.fields)-1 {
		w.current = 0
	}
	w.current++
	img := data.NewGradientFromField(w.fields[w.current], 720, 720)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.current+1, len(w.fields)))
}
func (w *PreviewWindow) Previous() {
	if w.current == 1 {
		w.current = len(w.fields)
	}
	w.current--
	img := data.NewGradientFromField(w.fields[w.current], 720, 720)
	w.image.SetImage(img)
	w.pageLb.Set(fmt.Sprintf("%d/%d", w.current+1, len(w.fields)))
}
func (w *PreviewWindow) PlayPause() {
	if w.pause {
		w.nextBtn.Disable()
		w.prevBtn.Disable()
		ticker := time.NewTicker(50 * time.Millisecond)
		w.timerQ = make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					if w.current >= len(w.fields) {
						w.current = 0
					}
					w.Next()
				case <-w.timerQ:
					ticker.Stop()
					return
				}
			}
		}()
		w.pause = false
	} else {
		w.nextBtn.Enable()
		w.prevBtn.Enable()
		close(w.timerQ)
		w.pause = true
	}
}
func (w *PreviewWindow) Show() { w.window.Show() }
func (w *PreviewWindow) Hide() { w.window.Hide() }
