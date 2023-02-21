package zap

import (
	"errors"
	"go.uber.org/zap"
	"testing"
)

func TestLogger(t *testing.T) {
	log := NewLogger(zap.NewNop())
	log.D("debug")
	log.I("info")
	log.W("warning")
	log.E(errors.New("error"))
	log.P("panic!")
}
