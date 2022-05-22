package keyring

import (
	"github.com/deifyed/xctl/pkg/tools/secrets"
)

const secretServiceItemNotFound = "The specified item could not be found in the keyring"

func handleError(err error, defaultError error) error {
	switch err.Error() {
	case secretServiceItemNotFound:
		return secrets.ErrNotFound
	default:
		return defaultError
	}
}
