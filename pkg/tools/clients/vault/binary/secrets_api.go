package binary

import (
	"fmt"
	"io"
	"strings"
)

func (c *client) Put(name string, secretAttributes map[string]string) error {
	_, err := c.runVaultCommand(
		"kv",
		"put",
		fmt.Sprintf("secret/%s", name),
		strings.Join(attributesAsArray(secretAttributes), " "),
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}

func (c *client) Get(name string, key string) (string, error) {
	result, err := c.runVaultCommand(
		"kv",
		"get",
		"-field", key,
		fmt.Sprintf("secret/%s", name),
	)
	if err != nil {
		return "", fmt.Errorf("running command: %w", err)
	}

	buf, err := io.ReadAll(result)
	if err != nil {
		return "", fmt.Errorf("buffering: %w", err)
	}

	return string(buf), nil
}

func (c *client) Delete(name string) error {
	_, err := c.runVaultCommand(
		"kv",
		"metadata",
		"delete",
		fmt.Sprintf("secret/%s", name),
	)
	if err != nil {
		return fmt.Errorf("running command: %w", err)
	}

	return nil
}
