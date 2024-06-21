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
	widget.BaseWidget
	fields       []*data.VelocityField
	parentWindow fyne.Window
	app          fyne.App
	progressBar  *widget.ProgressBar
	button       *widget.Button
	previewBtn   *widget.Button
	onSuccess    func(fields []*data.VelocityField)
}

func NewFileImporter(app fyne.App, window fyne.Window, name string, onSuccess func(fields []*data.VelocityField)) *FileImporter {
	fi := &FileImporter{
		fields:       make([]*data.VelocityField, 0),
		parentWindow: window,
		app:          app,
		onSuccess:    onSuccess,
	}
	fi.button = widget.NewButton(name, fi.onButton())
	fi.previewBtn = widget.NewButton("Просмотр", fi.onPreview())
	fi.progressBar = widget.NewProgressBar()
	fi.progressBar.Hide()
	fi.previewBtn.Disable()
	fi.ExtendBaseWidget(fi)
	return fi
}

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
			fi.progressBar.Show()
			sort.Sort(utils.ByNumericalFilename(readers))
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
				fi.fields[i] = data.NewVelocityField(d)
				fi.progressBar.SetValue(float64(i) / float64(len(readers)))
			}
			fi.onSuccess(fi.fields)
			fi.progressBar.Hide()
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

func (fi *FileImporter) CreateRenderer() fyne.WidgetRenderer {
	c := container.New(
		layout.NewGridLayoutWithRows(2),
		container.New(layout.NewHBoxLayout(), fi.button, fi.previewBtn),
		fi.progressBar,
	)
	return widget.NewSimpleRenderer(c)
}
