/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */
package logger

import (
	"fmt"
	"os"
)

type Logger struct {
	requestId string
}

func NewLogger(requestId string) Logger {
	return Logger{requestId: requestId}
}

// A huge error, execution can't continue
func (l *Logger) Critical(format string, a ...interface{}) {
	fmt.Printf("<2>["+l.requestId+"] "+format+"\n", a...)
	os.Exit(1)
}

// Not severe error
func (l *Logger) Error(format string, a ...interface{}) {
	fmt.Printf("<3>["+l.requestId+"] "+format+"\n", a...)
}

// An error might occur
func (l *Logger) Warning(format string, a ...interface{}) {
	fmt.Printf("<4>["+l.requestId+"] "+format+"\n", a...)
}

// Not an error, just unusual
func (l *Logger) Notice(format string, a ...interface{}) {
	fmt.Printf("<5>["+l.requestId+"] "+format+"\n", a...)
}

// Normal operations
func (l *Logger) Information(format string, a ...interface{}) {
	fmt.Printf("<6>["+l.requestId+"] "+format+"\n", a...)
}
