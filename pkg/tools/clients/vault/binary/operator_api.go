package binary

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"
)

func (c *client) Initialize() (vault.InitializationResponse, error) {
	result, err := c.runVaultCommand("operator", "init", "-format=json")
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("running command: %w", err)
	}

	response, err := parseInitializationResponse(result)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("parsing response: %w", err)
	}

	return response, nil
}

func (c *client) Unseal(key string) error {
	_, err := c.runVaultCommand("operator", "unseal", key)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}

func (c *client) EnableKv2() error {
	_, err := c.runVaultCommand("secrets", "enable", "kv-v2")
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
