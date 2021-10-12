package binary

import (
	"bytes"
	"fmt"
	"github.com/spf13/afero"
	"io"
	"os/exec"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
)

func (k kubectlBinaryClient) PodExec(opts kubectl.PodExecOpts) error {
	cmd := exec.Command(k.kubectlPath, fmt.Sprintf(
		"--namespace %s exec -it %s -- %s",
		opts.Pod.Namespace,
		opts.Pod.Name,
		opts.Command,
	))

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		k.logger.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return fmt.Errorf("executing pod command: %s", stderr.String())
	}

	if opts.Stdout == nil {
		return nil
	}

	_, err = io.Copy(cmd.Stdout, &stdout)
	if err != nil {
		return fmt.Errorf("pushing stdout data: %w", err)
	}

	return nil
}

func (k kubectlBinaryClient) PortForward(opts kubectl.PortForwardOpts) (kubectl.StopFn, error) {
	cmd := exec.Command(k.kubectlPath, fmt.Sprintf(
		"--namespace %s port-forward %s %d:%d",
		opts.Pod.Namespace,
		opts.Pod.Name,
		opts.PortFrom,
		opts.PortTo,
	))

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		k.logger.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return nil, fmt.Errorf("executing pod command: %s", err)
	}

	return func() error {
		return cmd.Process.Kill()
	}, nil
}

func NewKubectlBinaryClient(logger *logrus.Logger, fs *afero.Afero, kubeConfigPath string) (kubectl.Client, error) {
	kubectlPath, err := getKubectlPath(fs)
	if err != nil {
		return nil, fmt.Errorf("acquiring kubectl path: %w", err)
	}

	return &kubectlBinaryClient{
		logger:      logger,
		kubectlPath: kubectlPath,
		env: map[string]string{
			kubeConfigPathKey: kubeConfigPath,
		},
	}, nil
}
