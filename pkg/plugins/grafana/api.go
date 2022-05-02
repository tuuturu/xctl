package grafana

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	"github.com/deifyed/xctl/pkg/tools/secrets/kubernetes"
)

func PortForwardOpts() kubectl.PortForwardOpts {
	plugin := NewPlugin()

	return kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      plugin.Metadata.Name,
			Namespace: plugin.Metadata.Namespace,
		},
		ServicePort: grafanaPort,
		LocalPort:   grafanaLocalPort,
	}
}

func Credentials(client kubectl.Client) (v1alpha1.PluginCredentials, error) {
	secretClient := kubernetes.New(client, pluginNamespace)

	username, err := secretClient.Get(secretName(), adminUsernameKey)
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("retrieving username: %w", err)
	}

	password, err := secretClient.Get(secretName(), adminPasswordKey)
	if err != nil {
		return v1alpha1.PluginCredentials{}, fmt.Errorf("retrieving password: %w", err)
	}

	return v1alpha1.PluginCredentials{
		Username: username,
		Password: password,
	}, nil
}
