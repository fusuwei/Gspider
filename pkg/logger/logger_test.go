package logger

import (
	"fmt"
	"runtime"
	"testing"
)

func TestLogger(t *testing.T) {
	Debug("111111111")
}

func TestNewLogger(t *testing.T) {
	l := NewLogger("debug", "test1", false, "log", "test")
	l.Debugf("%s", "test")
	funcName, file, line, ok := runtime.Caller(10)
	if ok {
		fmt.Println("Func Name=" + runtime.FuncForPC(funcName).Name())
		fmt.Printf("file: %s    line=%d\n", file, line)
	}
}
