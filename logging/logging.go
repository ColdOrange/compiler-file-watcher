package logging

import (
	"io"
	"os"
	"time"

	"compiler-file-watcher/config"
	"github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

var levelMap = map[string]log.Level{
	"debug": log.DebugLevel,
	"info":  log.InfoLevel,
	"warn":  log.WarnLevel,
	"error": log.ErrorLevel,
	"fatal": log.FatalLevel,
}

func init() {
	// type
	logType := config.LoggingConfig.Type
	if logType == "text" {
		textFormatter := &log.TextFormatter{}
		textFormatter.FullTimestamp = true
		textFormatter.TimestampFormat = "2006-01-02 15:04:05.999999"
		log.SetFormatter(textFormatter)
	} else {
		jsonFormatter := &log.JSONFormatter{}
		jsonFormatter.TimestampFormat = "2006-01-02 15:04:05.999999"
		log.SetFormatter(jsonFormatter)
	}

	// level
	level, exist := levelMap[config.LoggingConfig.Level]
	if !exist {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}

	// console
	if config.LoggingConfig.Console {
		log.SetOutput(os.Stdout)
		return
	}

	// file
	logFileName := config.LoggingConfig.File

	// rotate
	rotate := config.LoggingConfig.Rotate

	logFile, err := rotatelogs.New(
		"./log/"+logFileName+".%Y%m%d%H",
		rotatelogs.WithRotationTime(time.Hour*time.Duration(rotate)),
	)

	var logWriter io.Writer
	if err != nil {
		logWriter = os.Stdout
	} else {
		logWriter = logFile
	}
	log.SetOutput(logWriter)
}

// Wrapper
func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Debugln(args ...interface{}) {
	log.Debugln(args...)
}

func Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Infoln(args ...interface{}) {
	log.Infoln(args...)
}

func Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Warnln(args ...interface{}) {
	log.Warnln(args...)
}

func Warnf(format string, args ...interface{}) {
	log.Warnf(format, args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Errorln(args ...interface{}) {
	log.Errorln(args...)
}

func Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

func Fatalln(args ...interface{}) {
	log.Fatalln(args...)
}

func Fatalf(format string, args ...interface{}) {
	log.Fatalf(format, args...)
}
