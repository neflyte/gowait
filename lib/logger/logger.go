package logger

import (
	"github.com/sirupsen/logrus"
)

const (
	LogFormatText = "text"
	LogFormatJSON = "json"

	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
	LogLevelFatal = "fatal"
)

var (
	rootLogger *logrus.Logger
)

func init() {
	rootLogger = logrus.New()
	rootLogger.SetLevel(logrus.DebugLevel)
	rootLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})
}

func ConfigureFormat(format string) {
	switch format {
	case LogFormatText:
		rootLogger.SetFormatter(&logrus.TextFormatter{
			ForceColors:      true,
			FullTimestamp:    true,
			QuoteEmptyFields: true,
		})
	case LogFormatJSON:
		rootLogger.SetFormatter(&logrus.JSONFormatter{})
	}
}

func ConfigureLevel(level string) {
	loggerLevel := logrus.InfoLevel
	switch level {
	case LogLevelDebug:
		loggerLevel = logrus.DebugLevel
	case LogLevelInfo:
		loggerLevel = logrus.InfoLevel
	case LogLevelWarn:
		loggerLevel = logrus.WarnLevel
	case LogLevelError:
		loggerLevel = logrus.ErrorLevel
	case LogLevelFatal:
		loggerLevel = logrus.FatalLevel
	}
	rootLogger.SetLevel(loggerLevel)
}

func GetRootLogger() logrus.FieldLogger {
	return rootLogger
}

func Function(funcName string) LogDelegate {
	return NewLogDelegate(rootLogger.WithField(FuncField, funcName))
}
