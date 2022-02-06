package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
	"github.com/deifyed/xctl/pkg/tools/logging"
)

func (k kubectlBinaryClient) PodExec(opts kubectl.PodExecOpts, args ...string) error {
	log := logging.GetLogger(logFeature, "podexec")

	staticArgs := []string{
		"exec",
		"-it",
		"--namespace", opts.Pod.Namespace,
		opts.Pod.Name,
		"--",
	}

	cmd := exec.Command(k.kubectlPath, append(staticArgs, args...)...) //nolint:gosec

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
	log := logging.GetLogger(logFeature, "portforward")

	cmd := exec.Command(k.kubectlPath, "port-forward", //nolint:gosec
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
		log.Debug("executing command", commandLogFields{
			Stdout: stdout.String(),
			Stderr: stderr.String(),
		})

		if isConnectionRefused(stderr.String()) {
			return nil, kubectl.ErrConnectionRefused
		}

		return nil, fmt.Errorf("executing pod command: %s", err)
	}

	time.Sleep(portForwardWaitSeconds * time.Second)

	return func() error {
		log.Debug("terminating port forward")

		return cmd.Process.Kill()
	}, nil
}
