package binary

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"sigs.k8s.io/yaml"
)

func (k kubectlBinaryClient) GetUserToken() (io.Reader, error) {
	var kubeConfig kubeConfig

	rawKubeConfig, err := os.ReadFile(k.env[kubeConfigPathKey])
	if err != nil {
		return nil, fmt.Errorf("reading: %w", err)
	}

	err = yaml.Unmarshal(rawKubeConfig, &kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling: %w", err)
	}

	ctx, err := getCurrentContext(kubeConfig, kubeConfig.CurrentContext)
	if err != nil {
		return nil, fmt.Errorf("getting current context: %w", err)
	}

	user, err := getUserForContext(kubeConfig, ctx)
	if err != nil {
		return nil, fmt.Errorf("getting current user: %w", err)
	}

	return bytes.NewReader([]byte(user.Token)), nil
}
