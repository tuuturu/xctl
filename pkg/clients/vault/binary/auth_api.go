package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/clients/vault"

	"github.com/sirupsen/logrus"
)

const defaultIssuer = "https://kubernetes.default.svc.cluster.local"

func (c *client) EnableKubernetesAuthentication() error {
	cmd := exec.Command(c.vaultPath, "auth", "enable", "kubernetes")

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

func (c *client) ConfigureKubernetesAuthentication(opts vault.ConfigureKubernetesAuthenticationOpts) error {
	args := []string{
		"write",
		"auth/kubernetes/config",
		fmt.Sprintf("kubernetes_host=%s", opts.Host.String()),
		fmt.Sprintf("token_reviewer_jwt=%s", opts.TokenReviewerJWT),
		fmt.Sprintf("kubernetes_ca_cert=%s", opts.CACert),
		fmt.Sprintf("issuer=%s", defaultIssuer),
	}

	cmd := exec.Command(c.vaultPath, args...)

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
