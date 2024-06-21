package widgets

import (
	"LPT/data"
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type SettingsWidget struct {
	widget.BaseWidget
	vVelocityImporter *FileImporter
	hVelocityImporter *FileImporter
	tStepInput        *ParsedFloatEntry
	tInterStepInput   *ParsedIntEntry
	traceChecker      *widget.Check
	sizeEntry         *SizeSettings
	params            data.FieldParams
	tracing           bool
	onParsed          func(params data.FieldParams, tracing bool, vField, hField []*data.VelocityField)
	onFailed          func(input string, err error)
}

func NewSettingsWidget(app fyne.App, window fyne.Window, onFailInput func(input string, err error), onParsed func(params data.FieldParams, tracing bool, vField, hField []*data.VelocityField)) *SettingsWidget {
	settings := &SettingsWidget{
		onParsed: onParsed,
		onFailed: onFailInput,
	}
	settings.sizeEntry = NewSizeSettings(data.Size{MinAxisX: -3, MaxAxisX: 3, MinAxisY: -3, MaxAxisY: 3}, func(size data.Size) {
		settings.params.Size = size
		settings.tryParse()
	})
	settings.vVelocityImporter = NewFileImporter(app, window, "Импортировать", func(fields []*data.VelocityField) {
		settings.tryParse()
	})
	settings.hVelocityImporter = NewFileImporter(app, window, "Импортировать", func(fields []*data.VelocityField) {
		settings.tryParse()
	})
	tStepParsed := func(value float64) {
		settings.params.TimeStep = value
		settings.tryParse()
	}
	settings.tStepInput = NewParsedFloatEntry(tStepParsed, onFailInput)
	tInterStepParsed := func(value int) {
		settings.params.InterStepCount = value
		settings.tryParse()
	}

	settings.traceChecker = widget.NewCheck("Показывать путь частиц", func(value bool) {
		settings.tracing = value
		settings.tryParse()
	})
	settings.tInterStepInput = NewParsedIntEntry(tInterStepParsed, onFailInput)
	settings.ExtendBaseWidget(settings)
	return settings
}

func (s *SettingsWidget) CreateRenderer() fyne.WidgetRenderer {
	gridVertical := layout.NewGridLayoutWithRows(2)
	title := canvas.NewText("Параметры", color.Black)
	title.TextSize = 30
	c := container.NewVBox(
		container.NewPadded(container.NewPadded(title)),
		container.New(layout.NewVBoxLayout(),
			container.NewPadded(
				container.New(gridVertical, widget.NewLabelWithStyle("Шаг интерполяции", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), s.tStepInput),
			),
			container.NewPadded(
				container.New(gridVertical, widget.NewLabelWithStyle("Интервал интерполяции", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}), s.tInterStepInput),
			),
		),
		s.traceChecker,
		container.NewPadded(
			container.NewVBox(
				widget.NewLabelWithStyle("Размер поля", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
				container.NewCenter(s.sizeEntry),
			),
		),
		container.NewVBox(
			widget.NewLabelWithStyle("Поля скорости", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
			container.NewCenter(s.vVelocityImporter),
			container.NewCenter(s.hVelocityImporter),
		),
	)
	return widget.NewSimpleRenderer(c)
}

func (s *SettingsWidget) tryParse() {
	if s.hVelocityImporter == nil || s.vVelocityImporter == nil {
		return
	}
	if s.hVelocityImporter.fields == nil || s.vVelocityImporter.fields == nil {
		s.onFailed("Указаны не все поля скорости", errors.New("velocity fields are required"))
		return
	}
	if len(s.hVelocityImporter.fields) == 0 || len(s.vVelocityImporter.fields) == 0 {
		s.onFailed("Неверный формат импортируемых полей скорости - не удалось считать данные", errors.New("invalid velocity fields import format"))
		return
	}
	if len(s.vVelocityImporter.fields) != len(s.hVelocityImporter.fields) {
		s.onFailed("Неверный формат импортируемых полей скорости - разное количество кадров", errors.New("velocity fields must be same length"))
		return
	}
	s.onParsed(s.params, s.tracing, s.vVelocityImporter.fields, s.hVelocityImporter.fields)
}
