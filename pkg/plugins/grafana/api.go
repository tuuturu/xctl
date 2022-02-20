package grafana

import "github.com/deifyed/xctl/pkg/tools/clients/kubectl"

func PortForwardOpts() kubectl.PortForwardOpts {
	plugin, _ := NewPlugin(NewPluginOpts{})

	return kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      plugin.Metadata.Name,
			Namespace: plugin.Metadata.Namespace,
		},
		ServicePort: grafanaPort,
		LocalPort:   grafanaLocalPort,
	}
}
