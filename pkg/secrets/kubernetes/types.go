package kubernetes

import "github.com/deifyed/xctl/pkg/clients/kubectl"

const secretKind = "Secret"

type client struct {
	kubernetesClient kubectl.Client
	namespace        string
}
