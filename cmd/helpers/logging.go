package helpers

import (
	"os"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/sirupsen/logrus"
)

func InitializeLogging() {
	level := os.Getenv("XCTL_LOG_LEVEL")

	if level == "" {
		level = "info"
	}

	logging.SetLevel(level)
	logging.SetOut(os.Stdout)

	logrus.SetFormatter(&logrus.JSONFormatter{})
}
