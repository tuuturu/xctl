package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/sirupsen/logrus"
)

func (c *client) Initialize() (vault.InitializationResponse, error) {
	cmd := exec.Command(c.vaultPath, "operator", "init", "-format=json") //nolint:gosec

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

		return vault.InitializationResponse{}, errorHandler(err, fmt.Errorf("executing command: %w", err))
	}

	response, err := parseInitializationResponse(&stdout)
	if err != nil {
		return vault.InitializationResponse{}, fmt.Errorf("parsing initialization response: %w", err)
	}

	return response, nil
}

func (c *client) Unseal(key string) error {
	cmd := exec.Command(c.vaultPath, "operator", "unseal", key) //nolint:gosec

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

		return errorHandler(err, fmt.Errorf("executing command: %w", err))
	}

	return nil
}

func (c *client) EnableKv2() error {
	cmd := exec.Command(c.vaultPath, "secrets", "enable", "kv-v2") //nolint:gosec

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

		return errorHandler(err, fmt.Errorf("executing command: %w", err))
	}

	return nil
}
