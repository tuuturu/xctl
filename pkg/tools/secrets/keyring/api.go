package keyring

import (
	"encoding/json"
	"fmt"

	"github.com/99designs/keyring"
)

// Put knows how to store secrets in a keyring
func (c Client) Put(name string, secrets map[string]string) error {
	var ring, err = keyring.Open(keyring.Config{ServiceName: generateServiceName(c.EnvironmentName)})
	if err != nil {
		return fmt.Errorf("opening keyring: %w", err)
	}

	payload, err := json.Marshal(secrets)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	err = ring.Set(keyring.Item{
		Key:  name,
		Data: payload,
	})
	if err != nil {
		return fmt.Errorf("storing: %w", err)
	}

	return nil
}

// Get knows how to retrieve secrets from a keyring
func (c Client) Get(name string, key string) (string, error) {
	ring, err := keyring.Open(keyring.Config{ServiceName: generateServiceName(c.EnvironmentName)})
	if err != nil {
		return "", fmt.Errorf("opening keyring: %w", err)
	}

	item, err := ring.Get(name)
	if err != nil {
		return "", handleError(fmt.Errorf("retrieving secret: %w", err))
	}

	var content map[string]string

	err = json.Unmarshal(item.Data, &content)
	if err != nil {
		return "", fmt.Errorf("unmarshalling: %w", err)
	}

	return content[key], nil
}

// Delete knows how to remove secrets from a keyring
func (c Client) Delete(name string) error {
	ring, err := keyring.Open(keyring.Config{ServiceName: generateServiceName(c.EnvironmentName)})
	if err != nil {
		return fmt.Errorf("opening keyring: %w", err)
	}

	err = ring.Remove(name)
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}
