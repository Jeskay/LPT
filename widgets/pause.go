package widgets

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Pause struct {
	widget.BaseWidget
	icon      *widget.Icon
	circle    *canvas.Circle
	container *fyne.Container
	paused    bool
	onClick   func()

	fadeTimer *time.Timer
}

func NewPauseWidget(onClick func()) *Pause {
	w := &Pause{
		onClick: onClick,
		paused:  true,
		circle:  canvas.NewCircle(color.RGBA{0, 0, 0, 150}),
		icon:    widget.NewIcon(theme.MediaPlayIcon()),
	}
	w.circle.Resize(fyne.NewSize(100, 100))
	pd4 := theme.Padding() * 2
	w.icon.Resize(fyne.NewSize(w.circle.Size().Width/2, w.circle.Size().Height/2))
	c := container.NewWithoutLayout(w.circle, w.icon)
	c.Resize(w.circle.Size())
	w.container = container.NewCenter(c)

	w.circle.Move(fyne.NewPos(-w.circle.Size().Width/2+pd4, -w.circle.Size().Height/2+pd4))
	w.icon.Move(fyne.NewPos(-w.icon.Size().Width/2+pd4, -w.icon.Size().Height/2+pd4))
	w.fadeTimer = time.AfterFunc(2*time.Second, func() {
		w.circle.Hide()
		w.icon.Hide()
	})
	w.ExtendBaseWidget(w)
	return w
}

func (p *Pause) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(p.container)
}

func (p *Pause) Tapped(e *fyne.PointEvent) {
	if p.paused {
		p.icon.SetResource(theme.MediaPauseIcon())
	} else {
		p.icon.SetResource(theme.MediaPlayIcon())
	}
	p.paused = !p.paused
	p.onClick()
}
func (p *Pause) MouseIn(e *desktop.MouseEvent) {

}
func (p *Pause) MouseMoved(e *desktop.MouseEvent) {
	if p.fadeTimer != nil {
		p.circle.Show()
		p.icon.Show()
		p.fadeTimer.Reset(2 * time.Second)
	}
}

func (p *Pause) MouseOut() {

}
