package log

import "os"
import "fmt"
import "mback/config"

var LOG_LEVEL = config.GetConfig().LogLevel

func log(level string, msg string, params ...interface{}) {
	str := fmt.Sprintf(level+msg, params...)
	fmt.Println(str)
}

func Debug(msg string, params ...interface{}) {
	if LOG_LEVEL > 0 {
		return
	}
	log("DEBUG: ", msg, params...)
}

func Info(msg string, params ...interface{}) {
	if LOG_LEVEL > 1 {
		return
	}
	log("", msg, params...)
}

func Warn(msg string, params ...interface{}) {
	if LOG_LEVEL > 2 {
		return
	}
	log(" WARN: ", msg, params...)
}

func Fatal(msg string, params ...interface{}) {
	log("FATAL: ", msg, params...)
	os.Exit(1)
}
