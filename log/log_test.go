package log

import (
	"errors"
	"testing"
)

func TestLog(t *testing.T) {
	stdLog := NewLog()
	stdLog.D("debug")
	stdLog.I("info")
	stdLog.W("warning")
	stdLog.E(errors.New("error"))
	stdLog.P("panic!")
}

func TestFmt(t *testing.T) {
	fmtLog := NewFmt()
	fmtLog.D("debug")
	fmtLog.I("info")
	fmtLog.W("warning")
	fmtLog.E(errors.New("error"))
	fmtLog.P("panic!")
}
