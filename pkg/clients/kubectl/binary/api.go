package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/spf13/afero"

	"github.com/sirupsen/logrus"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
)

func (k kubectlBinaryClient) PodExec(opts kubectl.PodExecOpts) error {
	cmd := exec.Command(k.kubectlPath, "exec",
		"-it",
		"--namespace", opts.Pod.Namespace,
		opts.Pod.Name,
		"--",
		string(opts.Command),
	)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		logrus.WithFields(logrus.Fields{
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
	cmd := exec.Command(k.kubectlPath, "port-forward",
		"--namespace", opts.Pod.Namespace,
		opts.Pod.Name,
		fmt.Sprintf("%s:%s", strconv.Itoa(opts.PortFrom), strconv.Itoa(opts.PortTo)),
	)

	stderr := bytes.Buffer{}
	stdout := bytes.Buffer{}

	cmd.Env = k.envAsArray()
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"stdout": stdout.String(),
			"stderr": stderr.String(),
		}).Debug("executing command")

		return nil, fmt.Errorf("executing pod command: %s", err)
	}

	time.Sleep(portforwardWaitSeconds * time.Second)

	return func() error {
		return cmd.Process.Kill()
	}, nil
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
