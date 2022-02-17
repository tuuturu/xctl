package grafana

import (
	"fmt"

	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"
	vaultClient "github.com/deifyed/xctl/pkg/tools/clients/vault"
	vaultBinary "github.com/deifyed/xctl/pkg/tools/clients/vault/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/deifyed/xctl/pkg/plugins/vault"
	"github.com/spf13/afero"
)

func openVaultConnection(kubectlClient kubectl.Client) (kubectl.StopFn, error) {
	vaultPlugin := vault.NewVaultPlugin()

	stopFn, err := kubectlClient.PortForward(kubectl.PortForwardOpts{
		Service: kubectl.Service{
			Name:      vaultPlugin.Metadata.Name,
			Namespace: vaultPlugin.Metadata.Namespace,
		},
		ServicePort: vaultClient.DefaultPort,
		LocalPort:   vaultClient.DefaultPort,
	})
	if err != nil {
		return nil, fmt.Errorf("setting up vault port forward: %w", err)
	}

	return stopFn, nil
}

func prepareClients(fs *afero.Afero, cluster v1alpha1.Environment) (clientContainer, error) {
	kubeConfigPath, err := config.GetAbsoluteKubeconfigPath(cluster.Metadata.Name)
	if err != nil {
		return clientContainer{}, fmt.Errorf("acquiring kube config path: %w", err)
	}

	kubectlClient, err := kubectlBinary.New(fs, kubeConfigPath)
	if err != nil {
		return clientContainer{}, fmt.Errorf("acquiring Kubectl client: %w", err)
	}

	helmClient, err := helmBinary.New(fs, kubeConfigPath)
	if err != nil {
		return clientContainer{}, fmt.Errorf("acquiring Helm client: %w", err)
	}

	vc, err := vaultBinary.New(fs)
	if err != nil {
		return clientContainer{}, fmt.Errorf("acquiring Vault client: %w", err)
	}

	return clientContainer{
		kubectl: kubectlClient,
		helm:    helmClient,
		secrets: vc,
	}, nil
}
