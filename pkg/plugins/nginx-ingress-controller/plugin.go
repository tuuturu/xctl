package ingress

import "github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

const nginxIngressControllerPluginName = "ingress-nginx"

func NewNginxIngressControllerPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(nginxIngressControllerPluginName)

	plugin.Metadata.Name = nginxIngressControllerPluginName
	plugin.Metadata.Namespace = "kube-system"
	plugin.Spec.Helm.Chart = "ingress-nginx/ingress-nginx"

	return plugin
}
