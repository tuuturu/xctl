package binary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return errorHandler(err, fmt.Errorf("executing pod command: %s", stderr.String()))
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
		"--namespace", opts.Service.Namespace,
		fmt.Sprintf("service/%s", opts.Service.Name),
		fmt.Sprintf("%s:%s", strconv.Itoa(opts.LocalPort), strconv.Itoa(opts.ServicePort)),
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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return nil, errorHandler(err, fmt.Errorf("executing pod command: %s", err))
	}

	time.Sleep(portForwardWaitSeconds * time.Second)

	return func() error {
		log.Debug("terminating port forward")

		return cmd.Process.Kill()
	}, nil
}

func (k kubectlBinaryClient) PodReady(pod kubectl.Pod) (bool, error) {
	log := logging.GetLogger(logFeature, "podReady")

	args := []string{
		"--namespace", pod.Namespace,
		"get", "pod",
		pod.Name,
		"--output=json",
	}

	cmd := exec.Command(k.kubectlPath, args...) //nolint:gosec

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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return false, errorHandler(err, fmt.Errorf("running command: %w", err))
	}

	var result getPodResult

	err = json.Unmarshal(stdout.Bytes(), &result)
	if err != nil {
		return false, fmt.Errorf("unmarshalling response: %w", err)
	}

	for _, containerStatus := range result.Status.ContainerStatuses {
		if !containerStatus.Ready {
			return false, nil
		}
	}

	return true, nil
}
