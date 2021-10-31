package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
)

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

func (c *client) Put(path string, values map[string]string) error {
	return nil
}

func (c *client) Get(path string) (map[string]string, error) {
	return nil, nil
}
