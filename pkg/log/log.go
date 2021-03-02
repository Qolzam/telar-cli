package log

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
)

func Info(format string, a ...interface{}) {
	info := color.New(color.FgWhite, color.BgGreen).SprintFunc()
	fmt.Printf("%s ", info("[INFO] "))
	fmt.Printf(format, a...)
	fmt.Println()
}

func Error(format string, a ...interface{}) {
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("%s ", red("[Error]"))
	fmt.Printf(format, a...)
	fmt.Println()
}

func InfoStruct(a ...interface{}) {
	spew.Sdump(a...)
}
