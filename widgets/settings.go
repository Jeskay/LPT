package widgets

import (
	"LPT/data"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SettingsMenu struct {
	uImporter   *FileImporter
	vImporter   *FileImporter
	pCountInput *ParsedIntEntry
	tStepInput  *ParsedFloatEntry
	submitBtn   *widget.Button
}

func NewSettingsMenu(window fyne.Window, onChanged func(field *data.Field)) *SettingsMenu {
	sm := &SettingsMenu{
		uImporter: NewFileImporter(window, "Выбрать файл u"),
		vImporter: NewFileImporter(window, "Выбрать файл w"),
	}
	pCountParsed := func(value int) {
		fmt.Println(value)
	}
	pCountFailed := func(input string, err error) {
		fmt.Println(err)
	}
	sm.pCountInput = NewParsedIntEntry(pCountParsed, pCountFailed)
	tStepParsed := func(value float64) {
		sm.uImporter.SetTimeStep(value)
		sm.vImporter.SetTimeStep(value)
	}
	tStepFailed := func(input string, err error) {
		fmt.Println(err)
	}
	sm.tStepInput = NewParsedFloatEntry(tStepParsed, tStepFailed)
	sm.submitBtn = widget.NewButton("Применить", func() {
		step, particleCount, _, _ := sm.GetData()

		field := data.NewField(particleCount, data.Size{
			MinAxisX: -3,
			MaxAxisX: 3,
			MinAxisY: -3,
			MaxAxisY: 3,
		}, step)
		xField := sm.uImporter.GetFields()[0]
		yField := sm.vImporter.GetFields()[0]
		field.SetVelocity(xField.Data, yField.Data)
		onChanged(field)
	})
	return sm
}

func (sm *SettingsMenu) GetForm() *widget.Form {
	form := widget.NewForm(
		widget.NewFormItem("Шаг", sm.tStepInput.entry),
		widget.NewFormItem("Количество частиц", sm.pCountInput.entry),
		widget.NewFormItem("", sm.uImporter.button),
		widget.NewFormItem("", sm.vImporter.button),
		widget.NewFormItem("", sm.submitBtn),
	)
	return form
}

func (sm *SettingsMenu) GetData() (step float64, particleCount int, uField []*data.VelocityField, vField []*data.VelocityField) {
	step = sm.tStepInput.GetValue()
	particleCount = sm.pCountInput.GetValue()
	uField = sm.uImporter.GetFields()
	vField = sm.vImporter.GetFields()
	return
}
