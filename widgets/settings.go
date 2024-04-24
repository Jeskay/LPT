package widgets

import (
	"LPT/data"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type SettingsMenu struct {
	uImporter       *FileImporter
	vImporter       *FileImporter
	pCountInput     *ParsedIntEntry
	tStepInput      *ParsedFloatEntry
	tInterStepInput *ParsedIntEntry
	submitBtn       *widget.Button
}

func NewSettingsMenu(window fyne.Window, onChanged func(manager *data.FieldManager)) *SettingsMenu {
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
	tInterStepParsed := func(value int) {
		fmt.Println(value)
	}
	tInterStepFailed := func(input string, err error) {
		fmt.Println(err)
	}
	sm.tInterStepInput = NewParsedIntEntry(tInterStepParsed, tInterStepFailed)
	sm.submitBtn = widget.NewButton("Применить", func() {
		step, interStepCount, particleCount, uFields, vFields := sm.GetData()
		params := data.FieldParams{
			ParticleCount: particleCount,
			Size: data.Size{
				MinAxisX: -3,
				MaxAxisX: 3,
				MinAxisY: -3,
				MaxAxisY: 3,
			},
			TimeStep:       step,
			InterStepCount: interStepCount,
		}
		fieldManager, err := data.NewFieldManager(params, uFields, vFields)
		if err != nil {
			panic(err)
		}
		onChanged(fieldManager)
	})
	return sm
}

func (sm *SettingsMenu) GetForm() *widget.Form {
	form := widget.NewForm(
		widget.NewFormItem("Шаг", sm.tStepInput.entry),
		widget.NewFormItem("Количество шагов между снимками", sm.tInterStepInput.entry),
		widget.NewFormItem("Количество частиц", sm.pCountInput.entry),
		widget.NewFormItem("", sm.uImporter.button),
		widget.NewFormItem("", sm.vImporter.button),
		widget.NewFormItem("", sm.submitBtn),
	)
	return form
}

func (sm *SettingsMenu) GetData() (step float64, interStepCount, particleCount int, uField []*data.VelocityField, vField []*data.VelocityField) {
	step = sm.tStepInput.GetValue()
	interStepCount = sm.tInterStepInput.GetValue()
	particleCount = sm.pCountInput.GetValue()
	uField = sm.uImporter.GetFields()
	vField = sm.vImporter.GetFields()
	return
}
