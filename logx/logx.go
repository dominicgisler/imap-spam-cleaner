package logx

import (
	"fmt"
	"log"
)

const (
	TypeInfo  = "INFO"
	TypeDebug = "DEBUG"
	TypeWarn  = "WARN"
	TypeError = "ERROR"
)

var verbose bool

func Init(verboseMode bool) {
	verbose = verboseMode
}

func Info(v ...interface{}) {
	Log(TypeInfo, v...)
}

func Infof(format string, v ...interface{}) {
	Logf(TypeInfo, format, v...)
}

func Debug(v ...interface{}) {
	Log(TypeDebug, v...)
}

func Debugf(format string, v ...interface{}) {
	Logf(TypeDebug, format, v...)
}

func Warn(v ...interface{}) {
	Log(TypeWarn, v...)
}

func Warnf(format string, v ...interface{}) {
	Logf(TypeWarn, format, v...)
}

func Error(v ...interface{}) {
	Log(TypeError, v...)
}

func Errorf(format string, v ...interface{}) {
	Logf(TypeError, format, v...)
}

func Log(tp string, v ...interface{}) {
	if tp == TypeDebug && !verbose {
		return
	}
	log.Printf("[%s] %s\n", tp, fmt.Sprint(v...))
}

func Logf(tp, format string, v ...interface{}) {
	if tp == TypeDebug && !verbose {
		return
	}
	log.Printf("[%s] "+format+"\n", append([]interface{}{tp}, v...)...)
}
