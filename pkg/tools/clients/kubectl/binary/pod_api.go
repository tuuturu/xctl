package binary

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"time"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

	"github.com/deifyed/xctl/pkg/tools/logging"
)

func (k kubectlBinaryClient) PodExec(opts kubectl.PodExecOpts, args ...string) error {
	staticArgs := []string{"exec", "-it", opts.Pod.Name, "--"}

	stdout, err := k.runCommand(runCommandOpts{
		Log:       logging.GetLogger(logFeature, "podexec"),
		Namespace: opts.Pod.Namespace,
		Args:      append(staticArgs, args...),
	})
	if err != nil {
		return fmt.Errorf("executing: %w", err)
	}

	if opts.Stdout == nil {
		return nil
	}

	_, err = io.Copy(opts.Stdout, stdout)
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
