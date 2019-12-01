package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

var (
	Logger *logrus.Logger
)

func init() {
	Logger = logrus.New()
	Logger.SetLevel(logrus.DebugLevel)
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		FullTimestamp:    true,
		QuoteEmptyFields: true,
	})
	Logger.SetOutput(os.Stdout)
}

func WithField(fieldName, fieldValue string) *logrus.Entry {
	return Logger.WithField(fieldName, fieldValue)
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	lf := logrus.Fields{}
	for key, intf := range fields {
		lf[key] = intf
	}
	return Logger.WithFields(lf)
}

func AddField(logger *logrus.Entry, fieldName, fieldValue string) *logrus.Entry {
	if logger == nil {
		return WithField(fieldName, fieldValue)
	}
	return logger.WithField(fieldName, fieldValue)
}

func AddFields(logger *logrus.Entry, fields map[string]interface{}) *logrus.Entry {
	if logger == nil {
		return WithFields(fields)
	}
	lf := logrus.Fields{}
	for key, intf := range fields {
		lf[key] = intf
	}
	return logger.WithFields(lf)
}
