package grafana

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/tools/secrets/kubernetes"

	helmBinary "github.com/deifyed/xctl/pkg/tools/clients/helm/binary"
	kubectlBinary "github.com/deifyed/xctl/pkg/tools/clients/kubectl/binary"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/spf13/afero"
)

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

	secretsClient := kubernetes.New(kubectlClient, pluginNamespace)

	return clientContainer{
		kubectl: kubectlClient,
		helm:    helmClient,
		secrets: secretsClient,
	}, nil
}

func secretName() string {
	return fmt.Sprintf("xctl-%s", pluginName)
}
