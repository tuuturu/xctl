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

type NewPluginOpts struct {
	// Host defines the hostname Grafana should be available at
	Host string
	// AdminUsername defines the username of the admin user
	AdminUsername string
	// AdminPassword defines the password of the admin user
	AdminPassword string
}

type clientContainer struct {
	kubectl kubectl.Client
	helm    helm.Client
	secrets secrets.Client
}

const (
	grafanaPort      = 80
	grafanaLocalPort = 8000
)
