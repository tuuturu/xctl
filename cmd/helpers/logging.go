package helpers

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitializeLogging() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.JSONFormatter{})
}
