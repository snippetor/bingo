package command

import "github.com/fatih/color"

func print(format string, v ...interface{}) {
	color.Magenta(format, v...)
}

func printInfo(format string, v ...interface{}) {
	color.Blue(format, v...)
}

func printSuccess(format string, v ...interface{}) {
	color.Green(format, v...)
}

func printError(format string, v ...interface{}) {
	color.Red("ERROR: "+format, v...)
}
