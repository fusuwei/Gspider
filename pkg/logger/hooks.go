package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

type Caller struct {
}

func (c *Caller) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (c *Caller) Fire(entry *logrus.Entry) error {
	entry.Data["appName"] = "MyAppName"
	return nil
}

func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		if !strings.HasPrefix(file, "logrus") {
			break
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
