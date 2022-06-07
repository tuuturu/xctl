package keyring

import (
	"github.com/deifyed/xctl/pkg/tools/i18n"
	"github.com/deifyed/xctl/pkg/tools/secrets"
)

const (
	secretServiceItemNotFound = "The specified item could not be found in the keyring"
	secretServiceUserAborted  = "Cannot get secret of a locked object"
)

func handleError(err error, defaultError error) error {
	switch err.Error() {
	case secretServiceItemNotFound:
		return secrets.ErrNotFound
	case secretServiceUserAborted:
		return &i18n.HumanReadableError{
			Content: secrets.ErrUserAborted.Error(),
			Key:     "secrets/userAborted",
		}
	default:
		return defaultError
	}
}
