package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/clients/vault"
	"github.com/sirupsen/logrus"
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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		if isConnectionRefused(err) {
			return vault.InitializationResponse{}, vault.ErrConnectionRefused
		}

		return vault.InitializationResponse{}, fmt.Errorf("executing command: %w", err)
	}

	response, err := parseInitializationResponse(&stdout)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("parsing initialization response: %w", err)
	}

	return response, nil
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

func (c *client) EnableKv2() error {
	cmd := exec.Command(c.vaultPath, "secrets", "enable", "kv-v2")

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
