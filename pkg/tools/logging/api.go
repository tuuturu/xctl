package logging

import "github.com/sirupsen/logrus"

func CreateEntry(logger *logrus.Logger, feature, activity string) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"feature":  feature,
		"activity": activity,
	})
}
