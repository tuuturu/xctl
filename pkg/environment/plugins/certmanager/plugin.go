package certmanager

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func newPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(certManagerName)

	plugin.Metadata.Name = certManagerName
	plugin.Metadata.Namespace = "kube-system"

	// URL: https://artifacthub.io/packages/helm/cert-manager/cert-manager
	plugin.Spec.Helm.Chart = "cert-manager"
	plugin.Spec.Helm.Version = "1.8.0"
	plugin.Spec.Helm.Values = valuesTemplate

	plugin.Spec.Helm.Repository.Name = "jetstack"
	plugin.Spec.Helm.Repository.URL = "https://charts.jetstack.io"

	return plugin
}

const certManagerName = "cert-manager"

//go:embed values.yaml
var valuesTemplate string //nolint:gochecknoglobals
