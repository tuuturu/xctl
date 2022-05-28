package binary

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/spf13/afero"
)

func (k kubectlBinaryClient) Apply(manifest io.Reader) error {
	_, err := k.runCommand(runCommandOpts{
		Log:   logging.GetLogger(logFeature, "apply"),
		Stdin: manifest,
		Args:  []string{"apply", "-f", "-"},
	})
	if err != nil {
		return fmt.Errorf("applying: %w", err)
	}

	return nil
}

func (k kubectlBinaryClient) Delete(manifest io.Reader) error {
	_, err := k.runCommand(runCommandOpts{
		Log:   logging.GetLogger(logFeature, "delete"),
		Stdin: manifest,
		Args:  []string{"delete", "-f", "-"},
	})
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

func (k kubectlBinaryClient) Get(selector kubectl.Selector) (io.Reader, error) {
	stdout, err := k.runCommand(runCommandOpts{
		Log:       logging.GetLogger(logFeature, "get"),
		Namespace: selector.Namespace,
		Args:      []string{"get", selector.Kind, selector.Name, "--output", "yaml"},
	})
	if err != nil {
		return nil, fmt.Errorf("retrieving: %w", err)
	}

	return stdout, nil
}

func (k kubectlBinaryClient) DeleteResource(selector kubectl.Selector) error {
	_, err := k.runCommand(runCommandOpts{
		Log:       logging.GetLogger(logFeature, "delete"),
		Namespace: selector.Namespace,
		Args:      []string{"delete", selector.Kind, selector.Name},
	})
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

func (k kubectlBinaryClient) IsReady(selector kubectl.Selector) (bool, error) {
	responseStream, err := k.runCommand(runCommandOpts{
		Log:       logging.GetLogger(logFeature, "isReady"),
		Namespace: selector.Namespace,
		Args:      []string{"get", selector.Kind, selector.Name, "--output", "json"},
	})
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
