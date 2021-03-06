package ingress

import (
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

type nginxIngressController struct {
	cloudProvider cloud.Provider
}

type determineActionOpts struct {
	Ctx    reconciliation.Context
	Helm   helm.Client
	Plugin v1alpha1.Plugin
}
