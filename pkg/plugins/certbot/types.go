package certbot

import (
	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/logging"
	"github.com/deifyed/xctl/pkg/tools/reconciliation"
)

const logFeature = "plugin/certbot"

type certbotReconciler struct {
	cloudProvider cloud.Provider
}

type determineActionOpts struct {
	Ctx    reconciliation.Context
	Helm   helm.Client
	Plugin v1alpha1.Plugin
	Logger logging.Logger
}
