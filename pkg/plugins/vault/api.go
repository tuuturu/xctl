package vault

import "github.com/deifyed/xctl/pkg/tools/clients/kubectl"

func PortForwardOpts() kubectl.PortForwardOpts {
	plugin := NewPlugin()

	return kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      plugin.Metadata.Name,
			Namespace: plugin.Metadata.Namespace,
		},
		ServicePort: port,
		LocalPort:   port,
	}
}
