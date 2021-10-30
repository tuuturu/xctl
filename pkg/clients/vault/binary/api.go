package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/clients/vault"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func (c *client) Initialize() (vault.InitializationResponse, error) {
	cmd := exec.Command(c.vaultPath, "operator", "init", "-format=json")

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = c.envAsArray()
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return vault.InitializationResponse{}, fmt.Errorf("executing command: %w", err)
	}

	response, err := parseInitializationResponse(&stdout)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("parsing initialization response: %w", err)
	}

	return response, nil
}

func (c *client) SetToken(token string) {
	c.env["VAULT_TOKEN"] = token
}

func (c *client) Unseal(key string) error {
	cmd := exec.Command(c.vaultPath, "operator", "unseal", key)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = c.envAsArray()
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return fmt.Errorf("executing command: %w", err)
	}

	return nil
}

func (c *client) envAsArray() []string {
	env := make([]string, len(c.env))
	index := 0

	for key, value := range c.env {
		env[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return env
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
