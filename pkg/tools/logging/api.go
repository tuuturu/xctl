package logging

import (
	"os"
)

func GetLogger(feature, activity string) Logger {
	opts := loggerOpts{
		Out:   os.Stdout,
		Level: "debug",
	}

	return getLogrusWrapper(opts, feature, activity)
}
