package binary

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func (c *client) Put(name string, secretAttributes map[string]string) error {
	cmd := exec.Command(c.vaultPath, //nolint:gosec
		"kv",
		"put",
		fmt.Sprintf("secret/%s", name),
		strings.Join(attributesAsArray(secretAttributes), " "),
	)

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

func (c *client) Get(name string, key string) (string, error) {
	cmd := exec.Command(c.vaultPath, //nolint:gosec
		"kv",
		"get",
		"-field", key,
		fmt.Sprintf("secret/%s", name),
	)

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

		return "", fmt.Errorf("executing command: %w", err)
	}

	return stdout.String(), nil
}

func (c *client) Delete(name string) error {
	cmd := exec.Command(c.vaultPath, //nolint:gosec
		"kv",
		"metadata",
		"delete",
		fmt.Sprintf("secret/%s", name),
	)

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
