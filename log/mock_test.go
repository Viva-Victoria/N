package log

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_TestLog(t *testing.T) {
	var l Logger = LoggerMock{}
	_ = l
}

func TestLoggerMock_D(t *testing.T) {
	var m LoggerMock

	var called bool
	m.Debug = func(message string) {
		called = true
	}
	m.D("debug")
	require.True(t, called)

	m.Debug = nil
	m.T = t
	m.D("debug")

	m.T = nil
	m.D("debug")
}

func TestLoggerMock_I(t *testing.T) {
	var m LoggerMock

	var called bool
	m.Info = func(message string) {
		called = true
	}
	m.I("info")
	require.True(t, called)

	m.Info = nil
	m.T = t
	m.I("info")

	m.T = nil
	m.I("info")
}

func TestLoggerMock_W(t *testing.T) {
	var m LoggerMock

	var called bool
	m.Warn = func(message string) {
		called = true
	}
	m.W("warn")
	require.True(t, called)

	m.Warn = nil
	m.T = t
	m.W("warn")

	m.T = nil
	m.W("warn")
}

func TestLoggerMock_E(t *testing.T) {
	var m LoggerMock

	var called bool
	m.Error = func(err error) {
		called = true
	}

	var err = errors.New("mock")
	m.E(err)
	require.True(t, called)

	m.Error = nil

	m.T = nil
	m.E(err)
}

func TestLoggerMock_P(t *testing.T) {
	var m LoggerMock

	var called bool
	m.Panic = func(p any) {
		called = true
	}

	var err = errors.New("mock")
	m.P(err)
	require.True(t, called)
}
