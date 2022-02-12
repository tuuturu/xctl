package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
)

func (k kubectlBinaryClient) Apply(manifest io.Reader) error {
	log := logging.GetLogger(logFeature, "apply")

	cmd := exec.Command(k.kubectlPath, "apply", "-f", "-") //nolint:gosec

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = manifest

	err := cmd.Run()
	if err != nil {
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		return fmt.Errorf("executing pod command: %s", err)
	}

	return nil
}

func (k kubectlBinaryClient) Get(namespace string, resourceType string, name string) (io.Reader, error) {
	log := logging.GetLogger(logFeature, "get")

	cmd := exec.Command(k.kubectlPath,
		"get",
		"--namespace", namespace,
		"--output", "yaml",
		resourceType, name,
	)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		return nil, fmt.Errorf("executing pod command: %s", err)
	}

	return &stdout, nil
}

func New(fs *afero.Afero, kubeConfigPath string) (kubectl.Client, error) {
	kubectlPath, err := getKubectlPath(fs)
	if err != nil {
		return nil, fmt.Errorf("acquiring kubectl path: %w", err)
	}

	return &kubectlBinaryClient{
		kubectlPath: kubectlPath,
		env: map[string]string{
			kubeConfigPathKey: kubeConfigPath,
		},
	}, nil
}
