package widgets

import (
	"LPT/data"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type SettingsMenu struct {
	uImporter       *FileImporter
	vImporter       *FileImporter
	tStepInput      *ParsedFloatEntry
	tInterStepInput *ParsedIntEntry
	editBtn         *widget.Button
	setSizeBtn      *widget.Button
	submitBtn       *widget.Button
	field           *data.Field
	traceChecker    *widget.Check
	sizeEntry       *SizeSettings
	manager         *data.FieldManager
}

type SizeSettings struct {
	widget.BaseWidget
	minEntryX *ParsedFloatEntry
	maxEntryX *ParsedFloatEntry
	minEntryY *ParsedFloatEntry
	maxEntryY *ParsedFloatEntry
	size      data.Size
}

var defaultParams = data.FieldParams{
	Size: data.Size{
		MinAxisX: -3,
		MaxAxisX: 3,
		MinAxisY: -3,
		MaxAxisY: 3,
	},
}

func NewSizeSettings(initialParams data.Size) *SizeSettings {
	settings := &SizeSettings{size: initialParams}
	onFailed := func(input string, err error) {
		fmt.Println(err)
	}
	settings.minEntryX = NewParsedFloatEntry(func(value float64) { settings.size.MinAxisX = value }, onFailed)
	settings.maxEntryX = NewParsedFloatEntry(func(value float64) { settings.size.MaxAxisX = value }, onFailed)
	settings.minEntryY = NewParsedFloatEntry(func(value float64) { settings.size.MinAxisY = value }, onFailed)
	settings.maxEntryY = NewParsedFloatEntry(func(value float64) { settings.size.MaxAxisY = value }, onFailed)
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

func NewSettingsMenu(app fyne.App, window fyne.Window, onChanged func(manager *data.FieldManager)) *SettingsMenu {
	sm := &SettingsMenu{
		uImporter: NewFileImporter(app, window, "Выбрать файл u"),
		vImporter: NewFileImporter(app, window, "Выбрать файл w"),
		field:     data.NewEmptyField(defaultParams.Size),
		sizeEntry: NewSizeSettings(data.Size{MinAxisX: -3, MaxAxisX: 3, MinAxisY: -3, MaxAxisY: 3}),
	}
	tStepParsed := func(value float64) {
		sm.uImporter.SetTimeStep(value)
		sm.vImporter.SetTimeStep(value)
	}
	tStepFailed := func(input string, err error) {
		fmt.Println(err)
	}
	sm.tStepInput = NewParsedFloatEntry(tStepParsed, tStepFailed)
	tInterStepParsed := func(value int) {
		fmt.Println(value)
	}
	tInterStepFailed := func(input string, err error) {
		fmt.Println(err)
	}
	sm.traceChecker = widget.NewCheck("Показывать путь частиц", nil)
	sm.tInterStepInput = NewParsedIntEntry(tInterStepParsed, tInterStepFailed)
	sm.editBtn = widget.NewButton("Редактировать поле", func() {
		menu := NewEditMenu(app, "Field Editor", 1080, 720, sm.field)
		menu.window.Show()
	})
	sm.setSizeBtn = widget.NewButton("Установить размер", func() {
		sm.field = data.NewEmptyField(sm.sizeEntry.GetSize())
	})
	sm.submitBtn = widget.NewButton("Применить", func() {
		tracing, step, interStepCount, uFields, vFields := sm.GetData()
		params := data.FieldParams{
			Size:           sm.sizeEntry.GetSize(),
			TimeStep:       step,
			InterStepCount: interStepCount,
		}
		fieldManager, err := data.NewFieldManager(params, tracing, sm.field, uFields, vFields)
		if err != nil {
			panic(err)
		}
		sm.manager = fieldManager
		onChanged(sm.manager)
	})
	return sm
}

func (sm *SettingsMenu) GetForm() *widget.Form {
	form := widget.NewForm(
		widget.NewFormItem("Шаг", sm.tStepInput),
		widget.NewFormItem("Количество шагов между снимками", sm.tInterStepInput),
		widget.NewFormItem("", sm.uImporter),
		widget.NewFormItem("", sm.vImporter),
		widget.NewFormItem("", widget.NewSeparator()),
		widget.NewFormItem("", sm.sizeEntry),
		widget.NewFormItem("", sm.setSizeBtn),
		widget.NewFormItem("", widget.NewSeparator()),
		widget.NewFormItem("", sm.editBtn),
		widget.NewFormItem("", sm.traceChecker),
		widget.NewFormItem("", sm.submitBtn),
	)
	return form
}

func (sm *SettingsMenu) GetData() (tracing bool, step float64, interStepCount int, uField []*data.VelocityField, vField []*data.VelocityField) {
	step = sm.tStepInput.GetValue()
	interStepCount = sm.tInterStepInput.GetValue()
	uField = sm.uImporter.GetFields()
	vField = sm.vImporter.GetFields()
	tracing = sm.traceChecker.Checked
	return
}
