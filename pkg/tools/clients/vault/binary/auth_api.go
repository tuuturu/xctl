package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/tools/clients/vault"

	"github.com/sirupsen/logrus"
)

func (c *client) EnableKubernetesAuthentication() error {
	cmd := exec.Command(c.vaultPath, "auth", "enable", "kubernetes") //nolint:gosec

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

func (c *client) ConfigureKubernetesAuthentication(opts vault.ConfigureKubernetesAuthenticationOpts) error {
	args := []string{
		"write",
		"auth/kubernetes/config",
		fmt.Sprintf("kubernetes_host=%s", opts.Host.String()),
		fmt.Sprintf("token_reviewer_jwt=%s", opts.TokenReviewerJWT),
		fmt.Sprintf("kubernetes_ca_cert=%s", opts.CACert),
		fmt.Sprintf("issuer=%s", opts.Issuer.String()),
	}

	cmd := exec.Command(c.vaultPath, args...) //nolint:gosec

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
