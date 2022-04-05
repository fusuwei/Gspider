package logger

import (
	"github.com/fusuwei/gspider/tools"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"time"
)

var levelMap = map[string]logrus.Level{
	"error": logrus.ErrorLevel,
	"debug": logrus.DebugLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
}

//var logging *Logger

type Logger struct {
	Path     string
	Name     string
	Level    string
	terminal bool
	logger   *logrus.Logger
}

//func init() {
//	logging = NewLogger("debug", "test", true, "log")
//}

func NewLogger(level, name string, terminal bool, p ...string) *Logger {
	newLogger := &Logger{
		logger:   logrus.New(),
		terminal: terminal,
	}
	if name == "" {
		newLogger.Name = "log"
	} else {
		newLogger.Name = name
	}
	if path, ok := tools.MakeDir(p...); ok {
		newLogger.Path = path
	} else {
		log.Fatal("创建文件失败")
	}

	newLogger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			"time": "@time",
			"msg":  "~msg",
		},
	})
	newLogger.SetLevel(level)

	newLogger.SetStdout(newLogger.fileStdout())
	return newLogger
}

func (logger *Logger) SetFormatter(formatter logrus.Formatter) {
	if formatter != nil {
		logger.logger.SetFormatter(formatter)
	}
}

func (logger *Logger) SetLevel(level string) {
	if l, ok := levelMap[level]; ok {
		logger.Level = level
		logger.logger.SetLevel(l)
	} else {
		logger.Level = "debug"
		logger.logger.SetLevel(logrus.DebugLevel)
	}
}

func (logger *Logger) fileStdout() io.Writer {
	if logger.terminal {
		return os.Stdout
	}
	logf, err := rotatelogs.New(
		logger.Path+"/"+logger.Name+".%Y-%m-%d.log",
		rotatelogs.WithLinkName(logger.Path+logger.Name+"link.log"),
		rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationTime(time.Hour*12),
	)
	if err != nil {
		return os.Stdout
	}
	return logf
}

func (logger *Logger) SetStdout(writer io.Writer) {
	logger.logger.SetOutput(writer)
}

//func GetLogger() *Logger {
//	return logging
//}

func (logger *Logger) GetLogger() *logrus.Logger {
	return logger.logger
}

func (logger *Logger) Panic(msg string) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Panic(msg)
}

func (logger *Logger) Fatal(msg string) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Fatal(msg)
}

func (logger *Logger) Error(msg string) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Debug(msg)
}

func (logger *Logger) Info(msg string) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Info(msg)
}

func (logger *Logger) Debug(msg string) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Info(msg)
}

func (logger *Logger) Warn(msg string) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Warn(msg)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Errorf(format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Infof(format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Debugf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Warnf(format, args...)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Panicf(format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logger.WithFields(logrus.Fields{"name": logger.Name}).Fatalf(format, args...)
}

//func Panic(msg string) {
//	logging.Panic(msg)
//}

//func Fatal(msg string) {
//	logging.Fatal(msg)
//}
//
//func Error(msg string) {
//	logging.Error(msg)
//}
//
//func Info(msg string) {
//	logging.Info(msg)
//}
//
//func Debug(msg string) {
//	logging.Info(msg)
//}
//
//func Warn(msg string) {
//	logging.Warn(msg)
//}
//
//func Errorf(format string, args ...interface{}) {
//	logging.Errorf(format, args...)
//}
//
//func Infof(format string, args ...interface{}) {
//	logging.Infof(format, args...)
//}
//
//func Debugf(format string, args ...interface{}) {
//	logging.Debugf(format, args...)
//}
//
//func Warnf(format string, args ...interface{}) {
//	logging.Warnf(format, args...)
//}
//
//func Panicf(format string, args ...interface{}) {
//	logging.Panicf(format, args...)
//}
//
//func Fatalf(format string, args ...interface{}) {
//	logging.Fatalf(format, args...)
//}
