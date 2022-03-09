package logger

import (
	"fmt"
	"os"
)

// A huge error, execution can't continue
func Critical(format string, a ...interface{}) {
	fmt.Printf("<2>"+format+"\n", a...)
	os.Exit(1)
}

// Not severe error
func Error(format string, a ...interface{}) {
	fmt.Printf("<3>"+format+"\n", a...)
}

// An error might occur
func Warning(format string, a ...interface{}) {
	fmt.Printf("<4>"+format+"\n", a...)
}

// Not an error, just unusual
func Notice(format string, a ...interface{}) {
	fmt.Printf("<5>"+format+"\n", a...)
}

// Normal operations
func Information(format string, a ...interface{}) {
	fmt.Printf("<6>"+format+"\n", a...)
}
