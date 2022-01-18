package vault

import (
	_ "embed"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

const vaultPluginName = "vault"

func NewVaultPlugin() v1alpha1.Plugin {
	plugin := v1alpha1.NewPlugin(vaultPluginName)

	plugin.Metadata.Name = vaultPluginName
	plugin.Metadata.Namespace = "kube-system"

	plugin.Spec.Helm.Chart = "vault"
	plugin.Spec.Helm.Version = "0.17.1"
	plugin.Spec.Helm.Values = vaultValuesTemplate

	plugin.Spec.Helm.Repository.Name = "hashicorp"
	plugin.Spec.Helm.Repository.URL = "https://helm.releases.hashicorp.com"

	return plugin
}

//go:embed values.yaml
var vaultValuesTemplate string //nolint:gochecknoglobals
