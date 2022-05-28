package binary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/spf13/afero"
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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return errorHandler(err, fmt.Errorf("executing pod command: %s", err))
	}

	return nil
}

func (k kubectlBinaryClient) Delete(manifest io.Reader) error {
	log := logging.GetLogger(logFeature, "delete")

	cmd := exec.Command(k.kubectlPath, "delete", "-f", "-") //nolint:gosec

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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return errorHandler(err, fmt.Errorf("executing pod command: %s", err))
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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return nil, errorHandler(err, fmt.Errorf("executing pod command: %s", err))
	}

	return &stdout, nil
}

func (k kubectlBinaryClient) DeleteResource(namespace string, kind string, name string) error {
	log := logging.GetLogger(logFeature, "delete")

	cmd := exec.Command(k.kubectlPath,
		"delete",
		"--namespace", namespace,
		kind, name,
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

		err = fmt.Errorf("%s: %w", stderr.String(), err)

		return errorHandler(err, fmt.Errorf("executing pod command: %s", err))
	}

	return nil
}

func (k kubectlBinaryClient) IsReady(selector kubectl.Selector) (bool, error) {
	log := logging.GetLogger(logFeature, "isReady")

	responseStream, err := k.runCommand(log,
		"--namespace", selector.Namespace,
		"get", selector.Kind, selector.Name,
		"--output", "json",
	)
	if err != nil {
		return false, fmt.Errorf("querying: %w", err)
	}

	rawResponse, err := io.ReadAll(responseStream)
	if err != nil {
		return false, fmt.Errorf("buffering: %w", err)
	}

	var response struct {
		Status struct {
			AvailableReplicas int `json:"availableReplicas"`
		} `json:"status"`
	}

	err = json.Unmarshal(rawResponse, &response)
	if err != nil {
		return false, fmt.Errorf("unmarshalling: %w", err)
	}

	return response.Status.AvailableReplicas > 0, nil
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
