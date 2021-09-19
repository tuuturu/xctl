package preruns

import (
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/pkg/errors"
)

func ErrorTranslator(err error) string {
	switch {
	case errors.Is(err, config.ErrNotAuthenticated):
		return i18n.T(config.ErrNotAuthenticated.Error())
	default:
		return "could not find translation for error"
	}
}
