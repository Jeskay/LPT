package widgets

import (
	"strconv"

	"fyne.io/fyne/v2/widget"
)

type ParsedIntEntry struct {
	value    int
	entry    *widget.Entry
	onFailed func(input string, err error)
	onParsed func(value int)
}

type ParsedFloatEntry struct {
	value    float64
	entry    *widget.Entry
	onFailed func(input string, err error)
	onParsed func(value float64)
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

func NewParsedIntEntry(onParsed func(value int), onFailed func(input string, err error)) *ParsedIntEntry {
	pe := &ParsedIntEntry{
		entry:    widget.NewEntry(),
		onFailed: onFailed,
		onParsed: onParsed,
	}
	pe.entry.OnChanged = pe.parse
	return pe
}

func NewParsedFloatEntry(onParsed func(value float64), onFailed func(input string, err error)) *ParsedFloatEntry {
	pe := &ParsedFloatEntry{
		entry:    widget.NewEntry(),
		onFailed: onFailed,
		onParsed: onParsed,
	}
	pe.entry.OnChanged = pe.parse
	return pe
}
