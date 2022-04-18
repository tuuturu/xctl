package grafana

import (
	"fmt"

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

func Credentials(client kubectl.Client) (CredentialsContainer, error) {
	secretClient := kubernetes.New(client, pluginNamespace)

	username, err := secretClient.Get(secretName(), adminUsernameKey)
	if err != nil {
		return CredentialsContainer{}, fmt.Errorf("retrieving username: %w", err)
	}

	password, err := secretClient.Get(secretName(), adminPasswordKey)
	if err != nil {
		return CredentialsContainer{}, fmt.Errorf("retrieving password: %w", err)
	}

	return CredentialsContainer{
		Username: username,
		Password: password,
	}, nil
}
