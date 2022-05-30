package argocd

import (
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

type installer interface {
	Install(v1alpha1.Plugin) error
}

type applier interface {
	Apply(io.Reader) error
}

func installArgoCD(rctx reconciliation.Context, helmClient installer, kubectlClient applier, repo repository) error {
	plugin, err := NewPlugin()
	if err != nil {
		return fmt.Errorf("creating plugin: %w", err)
	}

	err = helmClient.Install(plugin)
	if err != nil {
		return fmt.Errorf("installing plugin: %w", err)
	}

	err = configureDeployKey(rctx, kubectlClient, repo)
	if err != nil {
		return fmt.Errorf("configuring deploy key: %w", err)
	}

	err = configureTrackedDirectories(rctx, kubectlClient)
	if err != nil {
		return fmt.Errorf("configuring tracked directories: %w", err)
	}

	return nil
}

func configureTrackedDirectories(rctx reconciliation.Context, kubectlClient applier) error {
	applicationsAppManifest, err := establishConfiguration(
		rctx.Filesystem,
		rctx.RootDirectory,
		rctx.EnvironmentManifest,
	)
	if err != nil {
		return fmt.Errorf("establishing configuration: %w", err)
	}

	err = kubectlClient.Apply(applicationsAppManifest)
	if err != nil {
		return fmt.Errorf("applying applications app manifest: %w", err)
	}

	namespacesAppManifest, err := establishNamespacesConfiguration(
		rctx.Filesystem,
		rctx.RootDirectory,
		rctx.EnvironmentManifest,
	)
	if err != nil {
		return fmt.Errorf("establishing namespaces configuration: %w", err)
	}

	err = kubectlClient.Apply(namespacesAppManifest)
	if err != nil {
		return fmt.Errorf("applying applications app manifest: %w", err)
	}
	return nil
}

func configureDeployKey(rctx reconciliation.Context, kubectlClient applier, repo repository) error {
	keys, err := generateKey()
	if err != nil {
		return fmt.Errorf("generating key pair: %w", err)
	}

	err = installDeployKey(rctx.Ctx, installDeployKeyOpts{
		SecretClient: rctx.Keyring,
		ClusterName:  rctx.EnvironmentManifest.Metadata.Name,
		Repository:   repo,
		PublicKey:    keys.PublicKey,
	})
	if err != nil {
		return fmt.Errorf("installing deploy key: %w", err)
	}

	secret, err := generateRepositorySecret(repo, keys.PrivateKey)
	if err != nil {
		return fmt.Errorf("generating repository secret: %w", err)
	}

	err = kubectlClient.Apply(secret)
	if err != nil {
		return fmt.Errorf("applying secret: %w", err)
	}

	return nil
}
