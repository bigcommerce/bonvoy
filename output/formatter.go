package output

import (
	"encoding/json"
	"fmt"
)

const (
	formatText = "text"
	formatJson = "json"
)

type Formatter interface {
	Format(o interface{}) (string, error)
}

func NewFormatter(format string) (formatter Formatter, err error)  {
	switch format {
	case formatText:
		formatter = NewTextFormatter()
	case formatJson:
		formatter = NewJsonFormatter()
	default:
		err = fmt.Errorf("unknown format: %s", format)
	}
	return formatter, err
}

// Text formatting
type TextFormatter struct {}

func NewTextFormatter() Formatter {
	return &TextFormatter{}
}

func (f *TextFormatter) Format(o interface{}) (string, error) {
	switch o.(type) {
	case string:
		return fmt.Sprint(o), nil
	default:
		return fmt.Sprintf("%v", o), nil
	}
}

// JSON formatter
type JsonFormatter struct {}

func NewJsonFormatter() Formatter {
	return &JsonFormatter{}
}

func (f *JsonFormatter) Format(o interface{}) (string, error) {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil { return "", err }

	return fmt.Sprint(string(b)), nil
}