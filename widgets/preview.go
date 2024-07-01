package widgets

import (
	"LPT/data"
	"image"

	"fyne.io/fyne/v2"
)

type PreviewWindow struct {
	display *DisplayMenu
	fields  []*data.VelocityField
	window  fyne.Window
}

func NewPreviewWindow(app fyne.App, title string, w, h float32, fields []*data.VelocityField) *PreviewWindow {
	preview := &PreviewWindow{
		window: app.NewWindow(title),
		fields: fields,
		display: NewDisplayMenuWidget(func(index int) image.Image {
			return data.NewGradientFromField(fields[index], 720, 720)
		}, len(fields)),
	}
	preview.window.SetContent(preview.display)
	return preview
}

func (w *PreviewWindow) Show() { w.window.Show() }
func (w *PreviewWindow) Hide() { w.window.Hide() }
