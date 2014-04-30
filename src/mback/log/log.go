package log

import (
	"fmt"
	"os"
)

var LogLevel uint8 = 1

func log(level string, msg string, params ...interface{}) {
	str := fmt.Sprintf(level+msg, params...)
	fmt.Println(str)
}

func Debug(msg string, params ...interface{}) {
	if LogLevel > 0 {
		return
	}
	log("DEBUG: ", msg, params...)
}

func Info(msg string, params ...interface{}) {
	if LogLevel > 1 {
		return
	}
	log("", msg, params...)
}

func Warn(msg string, params ...interface{}) {
	if LogLevel > 2 {
		return
	}
	log(" WARN: ", msg, params...)
}

func Fatal(msg string, params ...interface{}) {
	log("FATAL: ", msg, params...)
	os.Exit(1)
}
