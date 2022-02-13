package kubernetes

import (
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
)

const secretKind = "Secret"

type client struct {
	kubernetesClient kubectl.Client
	namespace        string
}
