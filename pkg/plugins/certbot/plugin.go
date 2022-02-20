package certbot

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func newCertbotPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(certbotPluginName)

	plugin.Metadata.Name = certbotPluginName
	plugin.Metadata.Namespace = "kube-system"

	// URL: https://artifacthub.io/packages/helm/cert-manager/cert-manager/1.7.1
	plugin.Spec.Helm.Chart = "cert-manager"
	plugin.Spec.Helm.Version = "1.7.1"
	plugin.Spec.Helm.Values = valuesTemplate

	plugin.Spec.Helm.Repository.Name = "jetstack"
	plugin.Spec.Helm.Repository.URL = "https://charts.jetstack.io"

	return plugin
}

const certbotPluginName = "cert-manager"

//go:embed values.yaml
var valuesTemplate string //nolint:gochecknoglobals
