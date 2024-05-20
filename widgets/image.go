package widgets

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type ImageDisplay struct {
	widget.BaseWidget
	img *canvas.Image
}

func NewImageDisplay(img image.Image) *ImageDisplay {
	imageDis := &ImageDisplay{
		img: canvas.NewImageFromImage(img),
	}
	imageDis.img.FillMode = canvas.ImageFillOriginal
	imageDis.ExtendBaseWidget(imageDis)
	return imageDis
}

func (imageDis *ImageDisplay) SetImage(img image.Image) {
	imageDis.img.Image = img
	imageDis.img.Refresh()
}

func (imageDis *ImageDisplay) GetCanvas() *canvas.Image { return imageDis.img }

func (imageDis *ImageDisplay) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(nil, nil, nil, imageDis.img)
	return widget.NewSimpleRenderer(c)
}
