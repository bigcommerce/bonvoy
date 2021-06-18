package commands

import (
	"bonvoy/output"
	"fmt"
	"github.com/fatih/color"
)

var green = color.New(color.FgGreen)

func Ok(s ...interface{}) string {
	return green.Sprintln(s...)
}

func Info(s ...interface{}) string {
	return fmt.Sprintln(s...)
}

func (r *Registry) Output(o interface{}) error {
	formatter, err := output.NewFormatter(r.GetOutputFormat())
	if err != nil { return err }

	o, oErr := formatter.Format(o)
	if oErr != nil { return oErr }

	fmt.Println(o)
	return nil
}