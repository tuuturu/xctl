package vault

import (
	"fmt"

	helmBinary "github.com/deifyed/xctl/pkg/clients/helm/binary"
	kubectlBinary "github.com/deifyed/xctl/pkg/clients/kubectl/binary"
	vaultBinary "github.com/deifyed/xctl/pkg/clients/vault/binary"
	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
	"github.com/deifyed/xctl/pkg/clients/vault"
)

func initializeVault(kubectlClient kubectl.Client, vaultClient vault.Client) error {
	podSelector := kubectl.Pod{Name: "vault-0", Namespace: "kube-system"}

	stopFn, err := kubectlClient.PortForward(kubectl.PortForwardOpts{
		Pod:      podSelector,
		PortFrom: vault.DefaultPort,
		PortTo:   vault.DefaultPort,
	})
	if err != nil {
		return fmt.Errorf("forwarding vault port: %w", err)
	}

	defer func() {
		_ = stopFn()
	}()

	initResponse, err := vaultClient.Initialize()
	if err != nil {
		return fmt.Errorf("running init: %w", err)
	}

	vaultClient.SetToken(initResponse.RootToken)

	for index := 0; index < 3; index++ {
		err = vaultClient.Unseal(initResponse.UnsealKeysB64[index])
		if err != nil {
			return fmt.Errorf("unsealing: %w", err)
		}
	}

	return nil
}

func prepareClients(fs *afero.Afero, kubeConfigPath string) (clientContainer, error) {
	helmClient, err := helmBinary.New(fs, kubeConfigPath)
	if err != nil {
		return clientContainer{}, fmt.Errorf("creating helm binary client: %w", err)
	}

	kubectlClient, err := kubectlBinary.New(fs, kubeConfigPath)
	if err != nil {
		return clientContainer{}, fmt.Errorf("creating kubectl binary client: %w", err)
	}

	vaultClient, err := vaultBinary.New(fs)
	if err != nil {
		return clientContainer{}, fmt.Errorf("creating vault binary client: %w", err)
	}

	return clientContainer{
		kubectl: kubectlClient,
		vault:   vaultClient,
		helm:    helmClient,
	}, nil
}
