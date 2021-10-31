package vault

import (
	"github.com/deifyed/xctl/pkg/clients/helm"
	"github.com/deifyed/xctl/pkg/clients/kubectl"
	"github.com/deifyed/xctl/pkg/clients/vault"
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/vault"

type vaultReconciler struct {
	cloudProvider cloud.Provider
}

type clientContainer struct {
	kubectl kubectl.Client
	vault   vault.Client
	helm    helm.Client
}
