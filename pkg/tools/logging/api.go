package logging

import (
	"io"
	"os"
)

var (
	level Level     = LevelInfo //nolint:gochecknoglobals
	out   io.Writer = os.Stdout //nolint:gochecknoglobals
)

func SetLevel(newLevel string) {
	level = Level(newLevel)
}

func GetLevel() Level {
	return level
}

func SetOut(r io.Writer) {
	out = r
}

func GetLogger(feature, activity string) Logger {
	opts := loggerOpts{
		Out:   out,
		Level: level,
	}

	return getLogrusWrapper(opts, feature, activity)
}
