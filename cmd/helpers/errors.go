package helpers

import (
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/pkg/errors"
)

func ErrorTranslator(err error) string {
	switch {
	case errors.Is(err, cloud.ErrNotAuthenticated):
		return i18n.T(cloud.ErrNotAuthenticated.Error())
	default:
		return "could not find translation for error"
	}
}
