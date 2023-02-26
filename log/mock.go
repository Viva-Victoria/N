package log

import (
	"log"
	"testing"
)

type LoggerMock struct {
	T     *testing.T
	Debug func(message string)
	Info  func(message string)
	Warn  func(message string)
	Error func(err error)
	Panic func(panicErr any)
}

func (t LoggerMock) D(message string) {
	if t.Debug != nil {
		t.Debug(message)
		return
	}
	if t.T != nil {
		t.T.Log(message)
		return
	}

	log.Println(message)
}

func (t LoggerMock) I(message string) {
	if t.Info != nil {
		t.Info(message)
		return
	}
	if t.T != nil {
		t.T.Log(message)
		return
	}

	log.Println(message)
}

func (t LoggerMock) W(message string) {
	if t.Warn != nil {
		t.Warn(message)
		return
	}
	if t.T != nil {
		t.T.Log(message)
		return
	}

	log.Println(message)
}

func (t LoggerMock) E(err error) {
	if t.Error != nil {
		t.Error(err)
		return
	}
	if t.T != nil {
		t.T.Error(err)
		return
	}

	log.Println(err)
}

func (t LoggerMock) P(panicErr any) {
	if t.Panic != nil {
		t.Panic(panicErr)
		return
	}
	if t.T != nil {
		t.T.Fatal(panicErr)
		return
	}

	panic(panicErr)
}
