package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

const LogFormatText = "text"
const LogFormatJSON = "json"

var (
	rootLogger *logrus.Logger
)

func init() {
	rootLogger = logrus.New()
	rootLogger.SetLevel(logrus.DebugLevel)
	rootLogger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})
	rootLogger.SetOutput(os.Stdout)
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

func WithField(fieldName, fieldValue string) logrus.FieldLogger {
	return rootLogger.WithField(fieldName, fieldValue)
}

func WithFields(fields map[string]interface{}) logrus.FieldLogger {
	return rootLogger.WithFields(fields)
}

func AddField(log logrus.FieldLogger, fieldName, fieldValue string) logrus.FieldLogger {
	if log == nil {
		return WithField(fieldName, fieldValue)
	}
	return log.WithField(fieldName, fieldValue)
}

func AddFields(log logrus.FieldLogger, fields map[string]interface{}) logrus.FieldLogger {
	if log == nil {
		return WithFields(fields)
	}
	lf := logrus.Fields{}
	for key, intf := range fields {
		lf[key] = intf
	}
	return log.WithFields(lf)
}
