package log

import (
	"fmt"
	"log"
)

type Logger interface {
	D(message string)
	I(message string)
	W(message string)
	E(err error)
	P(panicErr any)
}

type StdLogger struct {
	println func(a any)
}

func NewLog() StdLogger {
	return StdLogger{
		println: func(a any) {
			log.Println(a)
		},
	}
}

func NewFmt() StdLogger {
	return StdLogger{
		println: func(a any) {
			fmt.Println(a)
		},
	}
}

func (l StdLogger) D(message string) {
	l.println(message)
}

func (l StdLogger) I(message string) {
	l.println(message)
}

func (l StdLogger) W(message string) {
	l.println(message)
}

func (l StdLogger) E(err error) {
	l.println(err)
}

func (l StdLogger) P(panicErr any) {
	l.println(panicErr)
}
