package widgets

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type ImageDisplay struct {
	widget.BaseWidget
	img       *canvas.Image
	circle    *canvas.Circle
	onDrag    func(imageDisplay *ImageDisplay, p fyne.Position)
	isKeyDown bool
	painting  bool
	crclSize  fyne.Size
}

func NewImageDisplay(img image.Image, radius float32, onDrag func(imageDisplay *ImageDisplay, p fyne.Position)) *ImageDisplay {
	imageDis := &ImageDisplay{
		img:       canvas.NewImageFromImage(img),
		circle:    canvas.NewCircle(color.Transparent),
		onDrag:    onDrag,
		isKeyDown: false,
		painting:  onDrag != nil,
		crclSize:  fyne.NewSize(radius, radius),
	}
	imageDis.img.SetMinSize(fyne.NewSize(float32(img.Bounds().Dx()), float32(img.Bounds().Dy())))
	imageDis.img.Resize(imageDis.img.MinSize())
	imageDis.img.Refresh()
	imageDis.img.FillMode = canvas.ImageFillOriginal
	imageDis.circle.Resize(imageDis.crclSize)
	imageDis.circle.StrokeWidth = 2
	imageDis.circle.StrokeColor = color.Black
	if !imageDis.painting {
		imageDis.circle.Hide()
	}
	imageDis.ExtendBaseWidget(imageDis)
	return imageDis
}

func (imageDis *ImageDisplay) SetImage(img image.Image) {
	imageDis.img.Image = img
	imageDis.img.Refresh()
}

func (imageDis *ImageDisplay) HideCircle() {
	imageDis.painting = false
	imageDis.circle.Hide()
}

func (imageDis *ImageDisplay) SetRadius(value float32) {
	imageDis.crclSize = fyne.NewSize(value, value)
	imageDis.circle.Resize(imageDis.crclSize)
	imageDis.circle.Refresh()
}

func (imageDis *ImageDisplay) ShowCircle() {
	imageDis.painting = true
	imageDis.circle.Show()
}

func (imageDis *ImageDisplay) GetCanvas() *canvas.Image { return imageDis.img }

func (imageDis *ImageDisplay) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewWithoutLayout(imageDis.img, imageDis.circle)
	return widget.NewSimpleRenderer(c)
}

func (imageDis *ImageDisplay) MouseMoved(e *desktop.MouseEvent) {
	nX := (imageDis.circle.Position1.X - imageDis.circle.Position2.X) / 2
	nY := (imageDis.circle.Position1.Y - imageDis.circle.Position2.Y) / 2
	imageDis.circle.Move(fyne.Position{X: e.Position.X + nX, Y: e.Position.Y + nY})
	if imageDis.isKeyDown {
		imageDis.onDrag(imageDis, e.Position)
	}
}

func (imageDis *ImageDisplay) MouseIn(e *desktop.MouseEvent) {}

func (imageDis *ImageDisplay) MouseOut() {}

func (imageDis *ImageDisplay) MouseDown(e *desktop.MouseEvent) {
	if imageDis.onDrag != nil {
		imageDis.isKeyDown = true
		imageDis.onDrag(imageDis, e.Position)
	}
}

func (imageDis *ImageDisplay) MouseUp(e *desktop.MouseEvent) {
	imageDis.isKeyDown = false
}
