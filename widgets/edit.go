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
	manager        *data.FieldManager
}

func NewEditWidget(window fyne.Window, fieldManager *data.FieldManager, w, h float32) *EditWidget {
	edit := &EditWidget{
		painting:       false,
		particleAmount: binding.NewInt(),
		spawnRadius:    binding.NewInt(),
		imageWidth:     w,
		imageHeight:    h,
		manager:        fieldManager,
	}
	edit.particleAmount.Set(10)
	edit.spawnRadius.Set(30)
	edit.colorPicker = NewColorPickerWidget(window, func(color color.RGBA) {
		edit.paintColor = color
	})
	img := fieldManager.Field.Image(500, 500)
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
	menu := container.NewHBox(
		container.NewGridWithRows(
			3,
			container.NewVBox(
				widget.NewLabelWithStyle("Режим редактирования", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				modeGroup,
			),
			container.NewVBox(
				widget.NewLabelWithStyle("Частицы", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				container.NewCenter(
					container.NewVBox(
						widget.NewLabel("Количество частиц"),
						NewParsedIntEntry(edit.onAmountParsed, nil, 100),
						widget.NewLabel("Радиус выделения"),
						NewIntSlider(edit.onRadiusParsed, 5, 200),
						container.NewPadded(edit.colorPicker),
					),
				),
			),
			container.NewVBox(
				widget.NewButton("Случайная генерация", edit.onRandom),
				widget.NewButton("Генерация вдоль оси X", edit.onLinear),
				widget.NewButton("Очистить", edit.onClear),
			),
		),
		edit.image,
	)
	return widget.NewSimpleRenderer(menu)
}

func (edit *EditWidget) onDrag(display *ImageDisplay, p fyne.Position) {
	size := edit.manager.Field.Size
	rad, err := edit.spawnRadius.Get()
	if err != nil {
		rad = 30
	}
	x := (float64(p.X) * (size.MaxAxisX - size.MinAxisX) / float64(display.img.Image.Bounds().Dx())) + size.MinAxisX
	y := (float64(p.Y) * (size.MinAxisY - size.MaxAxisY) / float64(display.img.Image.Bounds().Dy())) - size.MinAxisY
	radius := float64(rad) * (size.MaxAxisX - size.MinAxisX) / float64(display.img.Image.Bounds().Dx())
	if edit.painting {
		edit.manager.Field.SetColor(x, y, edit.paintColor, radius)
	} else {
		pCount, err := edit.particleAmount.Get()
		if err != nil {
			pCount = 10
		}
		edit.manager.Field.AddParticles(pCount, x, y, radius)
	}
	edit.updateImage()
}

func (edit *EditWidget) updateImage() {
	img := edit.manager.Field.Image(500, 500)
	edit.image.SetImage(resize.Resize(uint(edit.imageWidth), uint(edit.imageHeight), img, resize.Bilinear))
}

func (edit *EditWidget) onRandom() {
	amount, err := edit.particleAmount.Get()
	if err != nil || amount == 0 {
		amount = 1000
	}
	size := edit.manager.Field.Size
	*(edit.manager.Field) = *data.NewRandomField(amount, size)
	edit.updateImage()
}
func (edit *EditWidget) onLinear() {
	amount, err := edit.particleAmount.Get()
	if err != nil {
		amount = 1000
	}
	size := edit.manager.Field.Size
	*edit.manager.Field = *data.NewLinearField(amount, size)
	edit.updateImage()
}

func (edit *EditWidget) onClear() {
	size := edit.manager.Field.Size
	*edit.manager.Field = *data.NewEmptyField(size)
	edit.updateImage()
}

func (edit *EditWidget) onAmountParsed(value int) {
	edit.particleAmount.Set(value)
}
func (edit *EditWidget) onRadiusParsed(value int) {
	edit.spawnRadius.Set(value)
	if edit.manager != nil {
		edit.image.SetRadius(float32(value))
	}
}
