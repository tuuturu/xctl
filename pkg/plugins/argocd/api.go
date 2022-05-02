package argocd

import (
	"encoding/base64"
	"fmt"
	"github.com/deifyed/xctl/pkg/config"
	"io"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	"sigs.k8s.io/yaml"
)

func PortForwardOpts() kubectl.PortForwardOpts {
	return kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      "argocd-server",
			Namespace: config.DefaultOperationsNamespace,
		},
		ServicePort: 80,
		LocalPort:   8000,
	}
}

func Credentials(client kubectl.Client) (v1alpha1.PluginCredentials, error) {
	plugin, err := NewPlugin()
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("preparing plugin: %w", err)
	}

	secret, err := client.Get(plugin.Metadata.Namespace, "secret", "argocd-initial-admin-secret")
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("retrieving secret: %w", err)
	}

	raw, err := io.ReadAll(secret)
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("buffering: %w", err)
	}

	var response credentialsResponse

	err = yaml.Unmarshal(raw, &response)
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("unmarshalling: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(response.Data.Password)
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("decoding: %w", err)
	}

	return v1alpha1.PluginCredentials{
		Username: "admin",
		Password: string(decoded),
	}, nil
}

type credentialsResponse struct {
	Data struct {
		Password string `json:"password"`
	} `json:"data"`
}
