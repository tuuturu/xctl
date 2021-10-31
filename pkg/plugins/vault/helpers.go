package vault

import (
	"bytes"
	"fmt"
	"net/url"

	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/sirupsen/logrus"

	helmBinary "github.com/deifyed/xctl/pkg/clients/helm/binary"
	kubectlBinary "github.com/deifyed/xctl/pkg/clients/kubectl/binary"
	vaultBinary "github.com/deifyed/xctl/pkg/clients/vault/binary"
	"github.com/spf13/afero"

	"github.com/deifyed/xctl/pkg/clients/kubectl"
	"github.com/deifyed/xctl/pkg/clients/vault"
)

func installVault(clients clientContainer) error {
	log := logging.CreateEntry(logrus.StandardLogger(), logFeature, "installing")

	log.Debug("installing Helm chart")

	err := clients.helm.Install(NewVaultPlugin())
	if err != nil {
		return fmt.Errorf("installing Helm chart: %w", err)
	}

	podSelector := kubectl.Pod{Name: "vault-0", Namespace: "kube-system"}

	log.Debug("port forwarding vault container")

	stopFn, err := clients.kubectl.PortForward(kubectl.PortForwardOpts{
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

	log.Debug("initializing Vault")

	err = initializeVault(clients.vault)
	if err != nil {
		return fmt.Errorf("initializing vault: %w", err)
	}

	log.Debug("activating Kubernetes engine")

	err = activateKubernetesEngine(clients.vault, clients.kubectl, podSelector)
	if err != nil {
		return fmt.Errorf("activating Kubernetes engine: %w", err)
	}

	log.Debug("enabling kv-v2")

	err = clients.vault.EnableKv2()
	if err != nil {
		return fmt.Errorf("enabling kv-v2: %w", err)
	}

	return nil
}

func initializeVault(vaultClient vault.Client) error {
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

func activateKubernetesEngine(vaultClient vault.Client, kubectlClient kubectl.Client, podSelector kubectl.Pod) error {
	err := vaultClient.EnableKubernetesAuthentication()
	if err != nil {
		return fmt.Errorf("enabling Kubernetes authentication: %w", err)
	}

	host, err := acquireKubernetesHost(kubectlClient, podSelector)
	if err != nil {
		return fmt.Errorf("acquiring Kubernetes host: %w", err)
	}

	serviceToken, err := acquireServiceUserToken(kubectlClient, podSelector)
	if err != nil {
		return fmt.Errorf("acquiring service user token: %w", err)
	}

	caCert, err := acquireKubernetesCACertificate(kubectlClient, podSelector)
	if err != nil {
		return fmt.Errorf("acquiring Kubernetes CA certificate: %w", err)
	}

	err = vaultClient.ConfigureKubernetesAuthentication(vault.ConfigureKubernetesAuthenticationOpts{
		Host:             host,
		TokenReviewerJWT: serviceToken,
		CACert:           caCert,
		Issuer:           kubectl.DefaultKubernetesIssuer,
	})
	if err != nil {
		return fmt.Errorf("configuring Kubernetes authentication: %w", err)
	}

	return nil
}

func acquireKubernetesHost(client kubectl.Client, podSelector kubectl.Pod) (url.URL, error) {
	buf := bytes.Buffer{}

	err := client.PodExec(kubectl.PodExecOpts{
		Pod:    podSelector,
		Stdout: &buf,
	}, "printenv", "KUBERNETES_PORT_443_TCP_ADDR")
	if err != nil {
		return url.URL{}, fmt.Errorf("executing pod command: %w", err)
	}

	return url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s:443", buf.String()),
	}, nil
}

func acquireServiceUserToken(client kubectl.Client, podSelector kubectl.Pod) (string, error) {
	buf := bytes.Buffer{}

	err := client.PodExec(kubectl.PodExecOpts{
		Pod:    podSelector,
		Stdout: &buf,
	}, "cat", "/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return "", fmt.Errorf("executing pod command: %w", err)
	}

	return buf.String(), nil
}

func acquireKubernetesCACertificate(client kubectl.Client, podSelector kubectl.Pod) (string, error) {
	buf := bytes.Buffer{}

	err := client.PodExec(kubectl.PodExecOpts{
		Pod:    podSelector,
		Stdout: &buf,
	}, "cat", "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt")
	if err != nil {
		return "", fmt.Errorf("executing pod command: %w", err)
	}

	return buf.String(), nil
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
