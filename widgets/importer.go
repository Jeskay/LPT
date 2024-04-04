package widgets

import (
	"fmt"
	"os"

	"LPT/data"
	"LPT/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type FileImporter struct {
	step         float64
	fields       []*data.VelocityField
	parentWindow fyne.Window
	button       *widget.Button
}

func NewFileImporter(window fyne.Window, name string) *FileImporter {
	fi := &FileImporter{
		fields:       make([]*data.VelocityField, 0),
		parentWindow: window,
	}
	fi.button = widget.NewButton(name, fi.onButton())
	fi.button.Disable()
	return fi
}

func (fi *FileImporter) SetTimeStep(step float64) {
	fi.step = step
	if fi.button.Disabled() {
		fi.button.Enable()
	}
}

func (fi *FileImporter) GetWidget() *widget.Button        { return fi.button }
func (fi *FileImporter) GetFields() []*data.VelocityField { return fi.fields }

func (fi *FileImporter) onButton() func() {
	return func() {
		dialog.ShowFolderOpen(func(lUri fyne.ListableURI, err error) {
			if err != nil {
				fmt.Println(err)
				return
			}
			readers, err := lUri.List()
			if err != nil {
				fmt.Println(err)
				return
			}
			var currentTime float64 = 0
			fi.fields = make([]*data.VelocityField, len(readers))
			for i, uc := range readers {
				fmt.Println(uc.Path(), "\n", uc.Query())
				f, err := os.Open(uc.Path())
				if err != nil {
					fmt.Println(err)
					continue
				}
				d, err := utils.ReadDataFromFile(f)
				if err != nil {
					fmt.Println(err)
					continue
				}
				fi.fields[i] = data.NewVelocityField(d, currentTime)
				currentTime += fi.step
			}
		}, fi.parentWindow)
	}
}
