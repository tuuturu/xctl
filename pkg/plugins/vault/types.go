package vault

import (
	"github.com/deifyed/xctl/pkg/clients/helm"
	"github.com/deifyed/xctl/pkg/clients/kubectl"
	"github.com/deifyed/xctl/pkg/clients/vault"
)

type clientContainer struct {
	kubectl kubectl.Client
	vault   vault.Client
	helm    helm.Client
}
