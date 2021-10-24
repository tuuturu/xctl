package helpers

import (
	"github.com/sirupsen/logrus"
	"os"
)

func InitializeLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
