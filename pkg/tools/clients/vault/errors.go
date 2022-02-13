package vault

import "errors"

// ErrConnectionRefused represents inability to open a connection to Vault
var ErrConnectionRefused = errors.New("connection refused")
