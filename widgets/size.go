package widgets

import (
	"LPT/data"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SizeSettings struct {
	widget.BaseWidget
	minEntryX *ParsedFloatEntry
	maxEntryX *ParsedFloatEntry
	minEntryY *ParsedFloatEntry
	maxEntryY *ParsedFloatEntry
	size      data.Size
}

func NewSizeSettings(initialParams data.Size, onParsed func(data.Size)) *SizeSettings {
	settings := &SizeSettings{size: initialParams}
	onFailed := func(input string, err error) {
		fmt.Println(err)
	}
	settings.minEntryX = NewParsedFloatEntry(
		func(value float64) {
			settings.size.MinAxisX = value
			onParsed(settings.size)
		},
		onFailed,
		initialParams.MinAxisX,
	)
	settings.maxEntryX = NewParsedFloatEntry(
		func(value float64) {
			settings.size.MaxAxisX = value
			onParsed(settings.size)
		},
		onFailed,
		initialParams.MaxAxisX,
	)
	settings.minEntryY = NewParsedFloatEntry(
		func(value float64) {
			settings.size.MinAxisY = value
			onParsed(settings.size)
		},
		onFailed,
		initialParams.MinAxisY,
	)
	settings.maxEntryY = NewParsedFloatEntry(
		func(value float64) {
			settings.size.MaxAxisY = value
			onParsed(settings.size)
		},
		onFailed,
		initialParams.MaxAxisY,
	)
	settings.ExtendBaseWidget(settings)
	return settings
}

func (size *SizeSettings) GetSize() data.Size { return size.size }

func (size *SizeSettings) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewVBox(
		container.NewHBox(widget.NewLabel("Мин. оси абсцисс"), size.minEntryX),
		container.NewHBox(widget.NewLabel("Макс. оси абсцисс"), size.maxEntryX),
		container.NewHBox(widget.NewLabel("Мин. оси ординат"), size.minEntryY),
		container.NewHBox(widget.NewLabel("Макс. оси ординат"), size.maxEntryY),
	)
	return widget.NewSimpleRenderer(c)
}
