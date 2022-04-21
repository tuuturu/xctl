package kubernetes

import (
	"bytes"
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/tools/secrets"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"

	v1 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func (c client) Put(name string, secrets map[string]string) error {
	encoded := make(map[string][]byte)

	for key, value := range secrets {
		encoded[key] = []byte(value)
	}

	secret := v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       secretKind,
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      name,
			Namespace: c.namespace,
		},
		Type: v1.SecretTypeOpaque,
		Data: encoded,
	}

	rawManifest, err := yaml.Marshal(secret)
	if err != nil {
		return fmt.Errorf("marshalling manifest: %w", err)
	}

	err = c.kubernetesClient.Apply(bytes.NewReader(rawManifest))
	if err != nil {
		return fmt.Errorf("applying manifest: %w", err)
	}

	return nil
}

func (c client) Get(name string, key string) (string, error) {
	manifest, err := c.kubernetesClient.Get(c.namespace, secretKind, name)
	if err != nil {
		return "", fmt.Errorf("retrieving secret: %w", err)
	}

	rawManifest, err := io.ReadAll(manifest)
	if err != nil {
		return "", fmt.Errorf("buffering manifest: %w", err)
	}

	var secret v1.Secret

	err = yaml.Unmarshal(rawManifest, &secret)
	if err != nil {
		return "", fmt.Errorf("unmarshalling manifest: %w", err)
	}

	for currentKey, value := range secret.Data {
		if currentKey != key {
			continue
		}

		return string(value), nil
	}

	return "", fmt.Errorf("key %s not found", key)
}

func (c client) Delete(name string) error {
	err := c.kubernetesClient.DeleteResource(c.namespace, secretKind, name)
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

func New(kubernetesClient kubectl.Client, namespace string) secrets.Client {
	return &client{
		kubernetesClient: kubernetesClient,
		namespace:        namespace,
	}
}
