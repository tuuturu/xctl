package binary

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/spf13/afero"
)

func (c *client) SetToken(token string) {
	c.env["VAULT_TOKEN"] = token
}

func New(fs *afero.Afero) (vault.Client, error) {
	vaultPath, err := getVaultPath(fs)
	if err != nil {
		return nil, fmt.Errorf("acquiring vault path: %w", err)
	}

	return &client{
		vaultPath: vaultPath,
		env: map[string]string{
			"VAULT_ADDR": "http://127.0.0.1:8200",
			"PATH":       "/usr/bin",
		},
	}, nil
}
