package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type logrusWrapper struct {
	logger *logrus.Entry
}

func (l logrusWrapper) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l logrusWrapper) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l logrusWrapper) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l logrusWrapper) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l logrusWrapper) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l logrusWrapper) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l logrusWrapper) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l logrusWrapper) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func getLogrusWrapper(opts loggerOpts, feature, activity string) Logger {
	logger := logrus.New()

	level, _ := logrus.ParseLevel(string(opts.Level))

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(opts.Out)

	return logrusWrapper{
		logger: logger.WithFields(logrus.Fields{
			"tag": fmt.Sprintf("%s/%s", feature, activity),
		}),
	}
}
