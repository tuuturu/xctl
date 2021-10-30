package vault

import (
	"bytes"
	"fmt"
	"net/url"

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

func activateKubernetesEngine(vaultClient vault.Client, kubectlClient kubectl.Client) error {
	podSelector := kubectl.Pod{Name: "vault-0", Namespace: "kube-system"}

	stopFn, err := kubectlClient.PortForward(kubectl.PortForwardOpts{
		Pod:      podSelector,
		PortFrom: vault.DefaultPort,
		PortTo:   vault.DefaultPort,
	})
	if err != nil {
		return fmt.Errorf("port forwarding vault: %w", err)
	}

	defer func() {
		_ = stopFn()
	}()

	err = vaultClient.EnableKubernetesAuthentication()
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
