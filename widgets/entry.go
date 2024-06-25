package widgets

import (
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type ParsedIntEntry struct {
	widget.BaseWidget
	value    int
	entry    *widget.Entry
	onFailed func(input string, err error)
	onParsed func(value int)
}

type ParsedFloatEntry struct {
	widget.BaseWidget
	value    float64
	entry    *widget.Entry
	onFailed func(input string, err error)
	onParsed func(value float64)
}

type IntSlider struct {
	widget.BaseWidget
	min      int
	max      int
	value    binding.Float
	strValue binding.String
	onChange func(int)
}

func (pe *ParsedIntEntry) parse(input string) {
	v, err := strconv.Atoi(pe.entry.Text)
	if err != nil {
		pe.onFailed(pe.entry.Text, err)
		return
	}
	pe.value = v
	pe.onParsed(v)
}

func (pe *ParsedIntEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(pe.entry)
}

func (pe *ParsedFloatEntry) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(pe.entry)
}

func (s *IntSlider) CreateRenderer() fyne.WidgetRenderer {
	slider := widget.NewSliderWithData(float64(s.min), float64(s.max), s.value)
	var c *fyne.Container
	if s.strValue != nil {
		c = container.NewVBox(slider, container.NewCenter(widget.NewLabelWithData(s.strValue)))
	} else {
		c = container.NewPadded(slider)
	}
	return widget.NewSimpleRenderer(c)
}

func (pe *ParsedFloatEntry) parse(input string) {
	v, err := strconv.ParseFloat(pe.entry.Text, 64)
	if err != nil {
		pe.onFailed(pe.entry.Text, err)
		return
	}
	pe.value = v
	pe.onParsed(v)
}

func (pe *ParsedIntEntry) GetValue() int       { return pe.value }
func (pe *ParsedFloatEntry) GetValue() float64 { return pe.value }

func NewParsedIntEntry(onParsed func(value int), onFailed func(input string, err error), initialValue ...int) *ParsedIntEntry {
	pe := &ParsedIntEntry{
		entry:    widget.NewEntry(),
		onFailed: onFailed,
		onParsed: onParsed,
	}
	pe.entry.OnChanged = pe.parse
	if len(initialValue) > 0 {
		pe.entry.SetText(strconv.Itoa(initialValue[0]))
	}
	pe.ExtendBaseWidget(pe)
	return pe
}

func NewParsedFloatEntry(onParsed func(value float64), onFailed func(input string, err error), initialValue ...float64) *ParsedFloatEntry {
	pe := &ParsedFloatEntry{
		entry:    widget.NewEntry(),
		onFailed: onFailed,
		onParsed: onParsed,
	}
	pe.entry.OnChanged = pe.parse
	if len(initialValue) > 0 {
		pe.entry.SetText(strconv.FormatFloat(initialValue[0], 'f', -1, 64))
	}
	pe.ExtendBaseWidget(pe)
	return pe
}

func NewIntSliderWithIndicator(onChange func(int), min int, max int) *IntSlider {
	s := &IntSlider{
		onChange: onChange,
		min:      min,
		max:      max,
		value:    binding.NewFloat(),
		strValue: binding.NewString(),
	}
	s.value.AddListener(binding.NewDataListener(func() {
		v, err := s.value.Get()
		if err != nil {
			return
		}
		iv := int(math.Round(v))
		s.strValue.Set(strconv.Itoa(iv))
		s.onChange(iv)
	}))
	s.value.Set(float64(min))
	s.ExtendBaseWidget(s)
	return s
}

func NewIntSlider(onChange func(int), min int, max int) *IntSlider {
	s := &IntSlider{
		onChange: onChange,
		min:      min,
		max:      max,
		value:    binding.NewFloat(),
	}
	s.value.AddListener(binding.NewDataListener(func() {
		v, err := s.value.Get()
		if err != nil {
			return
		}
		iv := int(math.Round(v))
		s.onChange(iv)
	}))
	s.value.Set(float64(min))
	s.ExtendBaseWidget(s)
	return s
}
