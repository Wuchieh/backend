package logger

import (
	"string_backend_0001/internal/conf"
	"testing"
)

func init() {
	conf.Conf = conf.GetDefaultConfig()
	Init()
}

func TestLogger(t *testing.T) {
	Debug("hello world")
	Info("hello world")
	Warn("hello world")
	Error("hello world")
}
