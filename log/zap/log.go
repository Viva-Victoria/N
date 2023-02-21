package zap

import (
	"fmt"
	"go.uber.org/zap"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger(log *zap.Logger) Logger {
	return Logger{
		log: log,
	}
}

func (l Logger) D(message string) {
	l.log.Debug(message)
}

func (l Logger) I(message string) {
	l.log.Info(message)
}

func (l Logger) W(message string) {
	l.log.Warn(message)
}

func (l Logger) E(err error) {
	l.log.Error(err.Error())
}

func (l Logger) P(panicErr any) {
	l.log.Error(fmt.Sprintf("%v", panicErr))
}
