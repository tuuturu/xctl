package helpers

import (
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/pkg/errors"
)

func ErrorTranslator(err error) string {
	var humanReadableError *i18n.HumanReadableError

	if errors.As(err, &humanReadableError) {
		return i18n.T(humanReadableError.Key)
	}

	return err.Error()
}
