package logger

import "github.com/sirupsen/logrus"

const (
	FuncField = "func"
)

type LogDelegate interface {
	logrus.FieldLogger
	Err(err error) LogDelegate
	Field(name string, value interface{}) LogDelegate
	Fields(fields logrus.Fields) LogDelegate
	Function(funcName string) LogDelegate
}

type logDelegateData struct {
	delegate logrus.FieldLogger
}

func NewLogDelegate(existing logrus.FieldLogger) LogDelegate {
	if existing == nil {
		return &logDelegateData{
			delegate: GetRootLogger(),
		}
	}
	return &logDelegateData{
		delegate: existing,
	}
}

func (l logDelegateData) Err(err error) LogDelegate {
	return NewLogDelegate(l.delegate.WithError(err))
}

func (l logDelegateData) Field(name string, value interface{}) LogDelegate {
	return NewLogDelegate(l.delegate.WithField(name, value))
}

func (l logDelegateData) Fields(fields logrus.Fields) LogDelegate {
	return NewLogDelegate(l.delegate.WithFields(fields))
}

func (l logDelegateData) Function(funcName string) LogDelegate {
	return NewLogDelegate(l.delegate.WithField(FuncField, funcName))
}

func (l logDelegateData) WithField(key string, value interface{}) *logrus.Entry {
	return l.delegate.WithField(key, value)
}

func (l logDelegateData) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.delegate.WithFields(fields)
}

func (l logDelegateData) WithError(err error) *logrus.Entry {
	return l.delegate.WithError(err)
}

func (l logDelegateData) Debugf(format string, args ...interface{}) {
	l.delegate.Debugf(format, args...)
}

func (l logDelegateData) Infof(format string, args ...interface{}) {
	l.delegate.Infof(format, args...)
}

func (l logDelegateData) Printf(format string, args ...interface{}) {
	l.delegate.Printf(format, args...)
}

func (l logDelegateData) Warnf(format string, args ...interface{}) {
	l.delegate.Warnf(format, args...)
}

func (l logDelegateData) Warningf(format string, args ...interface{}) {
	l.delegate.Warningf(format, args...)
}

func (l logDelegateData) Errorf(format string, args ...interface{}) {
	l.delegate.Errorf(format, args...)
}

func (l logDelegateData) Fatalf(format string, args ...interface{}) {
	l.delegate.Fatalf(format, args...)
}

func (l logDelegateData) Panicf(format string, args ...interface{}) {
	l.delegate.Panicf(format, args...)
}

func (l logDelegateData) Debug(args ...interface{}) {
	l.delegate.Debug(args...)
}

func (l logDelegateData) Info(args ...interface{}) {
	l.delegate.Info(args...)
}

func (l logDelegateData) Print(args ...interface{}) {
	l.delegate.Print(args...)
}

func (l logDelegateData) Warn(args ...interface{}) {
	l.delegate.Warn(args...)
}

func (l logDelegateData) Warning(args ...interface{}) {
	l.delegate.Warning(args...)
}

func (l logDelegateData) Error(args ...interface{}) {
	l.delegate.Error(args...)
}

func (l logDelegateData) Fatal(args ...interface{}) {
	l.delegate.Fatal(args...)
}

func (l logDelegateData) Panic(args ...interface{}) {
	l.delegate.Panic(args...)
}

func (l logDelegateData) Debugln(args ...interface{}) {
	l.delegate.Debugln(args...)
}

func (l logDelegateData) Infoln(args ...interface{}) {
	l.delegate.Infoln(args...)
}

func (l logDelegateData) Println(args ...interface{}) {
	l.delegate.Println(args...)
}

func (l logDelegateData) Warnln(args ...interface{}) {
	l.delegate.Warnln(args...)
}

func (l logDelegateData) Warningln(args ...interface{}) {
	l.delegate.Warningln(args...)
}

func (l logDelegateData) Errorln(args ...interface{}) {
	l.delegate.Errorln(args...)
}

func (l logDelegateData) Fatalln(args ...interface{}) {
	l.delegate.Fatalln(args...)
}

func (l logDelegateData) Panicln(args ...interface{}) {
	l.delegate.Panicln(args...)
}
