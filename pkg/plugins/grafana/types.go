package grafana

import (
	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/tools/clients/helm"
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
	"github.com/deifyed/xctl/pkg/tools/secrets"
)

const logFeature = "plugin/grafana"

type reconciler struct {
	cloudProvider cloud.Provider
}

type valuesOpts struct {
	// SecretName defines the name of the secret of where to store plugin secrets
	SecretName string
	// SecretUsernameKey defines the key where the admin username is available
	SecretUsernameKey string
	// SecretPasswordKey defines the key where the admin password is available
	SecretPasswordKey string
}

type clientContainer struct {
	kubectl kubectl.Client
	helm    helm.Client
	secrets secrets.Client
}

const (
	grafanaPort      = 80
	grafanaLocalPort = 8000
	adminUsernameKey = "admin-user"
	adminPasswordKey = "admin-password"
)
