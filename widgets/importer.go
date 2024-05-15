package widgets

import (
	"fmt"
	"os"
	"sort"

	"LPT/data"
	"LPT/utils"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type FileImporter struct {
	step         float64
	fields       []*data.VelocityField
	parentWindow fyne.Window
	app          fyne.App
	elem         *fyne.Container
	prgBar       *widget.ProgressBar
	button       *widget.Button
	previewBtn   *widget.Button
}

func NewFileImporter(app fyne.App, window fyne.Window, name string) *FileImporter {
	fi := &FileImporter{
		fields:       make([]*data.VelocityField, 0),
		parentWindow: window,
		app:          app,
	}
	fi.button = widget.NewButton(name, fi.onButton())
	fi.previewBtn = widget.NewButton("Просмотр", fi.onPreview())
	fi.prgBar = widget.NewProgressBar()
	fi.elem = container.New(
		layout.NewVBoxLayout(),
		container.New(layout.NewHBoxLayout(), fi.button, fi.previewBtn),
		fi.prgBar,
	)
	fi.prgBar.Hide()
	fi.button.Disable()
	fi.previewBtn.Disable()
	return fi
}

func (fi *FileImporter) SetTimeStep(step float64) {
	fi.step = step
	if fi.button.Disabled() {
		fi.button.Enable()
	}
}

func (fi *FileImporter) GetWidget() *fyne.Container       { return fi.elem }
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
			fi.prgBar.Show()
			sort.Sort(utils.ByNumericalFilename(readers))
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
				fi.prgBar.SetValue(float64(i) / float64(len(readers)))
			}
			fi.prgBar.Hide()
			fi.button.SetText(lUri.Name())
		}, fi.parentWindow)
		fi.previewBtn.Enable()
	}
}

func (fi *FileImporter) onPreview() func() {
	return func() {
		w := NewPreviewWindow(fi.app, "Поле скорости", 720, 720, fi.fields)
		w.Show()
	}
}
