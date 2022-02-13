package vault

import (
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	"github.com/deifyed/xctl/pkg/tools/clients/vault"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

const logFeature = "plugin/vault"

type vaultReconciler struct {
	cloudProvider cloud.Provider
}

type determineActionOpts struct {
	rctx       reconciliation.Context
	helmClient helm.Client
	plugin     v1alpha1.Plugin
	kubectl    kubectl.Client
}

type clientContainer struct {
	kubectl kubectl.Client
	vault   vault.Client
	helm    helm.Client
}
