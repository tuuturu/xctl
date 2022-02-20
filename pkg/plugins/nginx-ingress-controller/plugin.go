package ingress

import "github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

func NewNginxIngressControllerPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(nginxIngressControllerPluginName)

	plugin.Metadata.Name = nginxIngressControllerPluginName
	plugin.Metadata.Namespace = "kube-system"

	// URL: https://github.com/kubernetes/ingress-nginx/
	plugin.Spec.Helm.Chart = "ingress-nginx"
	plugin.Spec.Helm.Version = "4.0.17"

	plugin.Spec.Helm.Repository.Name = "ingress-nginx"
	plugin.Spec.Helm.Repository.URL = "https://kubernetes.github.io/ingress-nginx"

	return plugin
}

const nginxIngressControllerPluginName = "ingress-nginx"
