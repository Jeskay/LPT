package widgets

import (
	"LPT/data"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/nfnt/resize"
)

const (
	Paint = "Раскрашивание"
	Edit  = "Изменение"
)

type EditMenu struct {
	window fyne.Window
	edit   *EditWidget
}

type EditWidget struct {
	widget.BaseWidget
	painting       bool
	particleAmount binding.Int
	spawnRadius    binding.Int
	image          *ImageDisplay
	imageHeight    float32
	imageWidth     float32
	paintColor     color.RGBA
	colorPicker    *ColorPickerWidget
	field          *data.Field
}

func NewEditMenu(app fyne.App, title string, w, h float32, field *data.Field) *EditMenu {
	menu := &EditMenu{
		window: app.NewWindow(title),
	}
	menu.edit = NewEditWidget(menu.window, field, w, h)
	menu.window.SetContent(menu.edit)
	return menu
}

func NewEditWidget(window fyne.Window, field *data.Field, w, h float32) *EditWidget {
	edit := &EditWidget{
		painting:       false,
		particleAmount: binding.NewInt(),
		spawnRadius:    binding.NewInt(),
		field:          field,
		imageWidth:     w,
		imageHeight:    h,
	}
	img := field.Image(500, 500)
	edit.particleAmount.Set(10)
	edit.spawnRadius.Set(30)
	edit.colorPicker = NewColorPickerWidget(window, func(color color.RGBA) {
		edit.paintColor = color
	})
	edit.image = NewImageDisplay(resize.Resize(uint(w), uint(h), img, resize.Bilinear), 30, edit.onDrag)
	edit.ExtendBaseWidget(edit)
	return edit
}

func (edit *EditWidget) CreateRenderer() fyne.WidgetRenderer {
	onModeChanged := func(mode string) {
		edit.painting = mode == Paint
	}
	modeGroup := widget.NewRadioGroup([]string{Paint, Edit}, onModeChanged)
	modeGroup.SetSelected(Paint)
	c := container.NewHBox(
		edit.image,
		container.NewVBox(
			modeGroup,
			widget.NewButton("Случайная генерация", edit.onRandom),
			widget.NewButton("Генерация вдоль оси X", edit.onLinear),
			widget.NewButton("Очистить", edit.onClear),
			edit.colorPicker,
			container.NewVBox(widget.NewLabel("Количество частиц"), NewIntSlider(edit.onAmountParsed, 1, 10000)),
			container.NewVBox(widget.NewLabel("Радиус выделения"), NewIntSlider(edit.onRadiusParsed, 5, 500)),
		),
	)
	return widget.NewSimpleRenderer(c)
}

func (edit *EditWidget) onDrag(display *ImageDisplay, p fyne.Position) {
	size := edit.field.Size
	rad, err := edit.spawnRadius.Get()
	if err != nil {
		rad = 30
	}
	x := (float64(p.X) * (size.MaxAxisX - size.MinAxisX) / float64(display.img.Image.Bounds().Dx())) + size.MinAxisX
	y := (float64(p.Y) * (size.MinAxisY - size.MaxAxisY) / float64(display.img.Image.Bounds().Dy())) - size.MinAxisY
	radius := float64(rad) * (size.MaxAxisX - size.MinAxisX) / float64(display.img.Image.Bounds().Dx())
	if edit.painting {
		edit.field.SetColor(x, y, edit.paintColor, radius)
	} else {
		pCount, err := edit.particleAmount.Get()
		if err != nil {
			pCount = 10
		}
		edit.field.AddParticles(pCount, x, y, radius)
	}
	edit.updateImage()
}

func (edit *EditWidget) updateImage() {
	img := edit.field.Image(500, 500)
	edit.image.SetImage(resize.Resize(uint(edit.imageWidth), uint(edit.imageHeight), img, resize.Bilinear))
}

func (edit *EditWidget) onRandom() {
	amount, err := edit.particleAmount.Get()
	if err != nil || amount == 0 {
		amount = 1000
	}
	size := edit.field.Size
	*edit.field = *data.NewRandomField(amount, size)
	edit.updateImage()
}
func (edit *EditWidget) onLinear() {
	amount, err := edit.particleAmount.Get()
	if err != nil {
		amount = 1000
	}
	size := edit.field.Size
	*edit.field = *data.NewLinearField(amount, size)
	edit.updateImage()
}

func (edit *EditWidget) onClear() {
	size := edit.field.Size
	*edit.field = *data.NewEmptyField(size)
	edit.updateImage()
}

func (edit *EditWidget) onAmountParsed(value int) {
	edit.particleAmount.Set(value)
}
func (edit *EditWidget) onRadiusParsed(value int) {
	edit.spawnRadius.Set(value)
	edit.image.SetRadius(float32(value))
}
