package logx

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
		PadLevelText:  true,
	})
}

func SetLevel(level string) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.Errorf("Invalid log level: %s", level)
		return
	}
	logger.SetLevel(lvl)
}

func Info(v ...any) {
	logger.Info(v...)
}

func Infof(format string, v ...any) {
	logger.Infof(format, v...)
}

func Debug(v ...any) {
	logger.Debug(v...)
}

func Debugf(format string, v ...any) {
	logger.Debugf(format, v...)
}

func Warn(v ...any) {
	logger.Warn(v...)
}

func Warnf(format string, v ...any) {
	logger.Warnf(format, v...)
}

func Error(v ...any) {
	logger.Error(v...)
}

func Errorf(format string, v ...any) {
	logger.Errorf(format, v...)
}

func Panic(v ...any) {
	logger.Panic(v...)
}

func Panicf(format string, v ...any) {
	logger.Panicf(format, v...)
}
