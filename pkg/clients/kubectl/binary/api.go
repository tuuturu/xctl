package binary

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
)

func (k kubectlBinaryClient) PodExec(opts kubectl.PodExecOpts) error {
	cmd := exec.Command(k.kubectlPath, fmt.Sprintf(
		"--namespace %s exec -it %s -- %s",
		opts.Pod.Namespace,
		opts.Pod.Name,
		opts.Command,
	))

	cmd.Env = []string{fmt.Sprintf("KUBECONFIG=%s", k.kubeConfigPath)}

	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	}

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("executing pod command: %s", stderr.String())
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

	cmd.Env = []string{fmt.Sprintf("KUBECONFIG=%s", k.kubeConfigPath)}

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("executing pod command: %s", stderr.String())
	}

	return func() error {
		return cmd.Process.Kill()
	}, nil
}

func NewKubectlBinaryClient(kubectlPath, kubeConfigPath string) kubectl.Client {
	return &kubectlBinaryClient{
		kubectlPath:    kubectlPath,
		kubeConfigPath: kubeConfigPath,
	}
}
