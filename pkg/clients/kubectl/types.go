package kubectl

import (
	"io"
	"net/url"
)

// DefaultKubernetesIssuer is the default issuer for Kubernetes found in .well-known
var DefaultKubernetesIssuer = url.URL{ //nolint:gochecknoglobals
	Scheme: "https",
	Host:   "kubernetes.default.svc.cluster.local",
}

type PodExecOpts struct {
	Pod    Pod
	Stdout io.Writer
}

type PortForwardOpts struct {
	Pod Pod

	PortFrom int
	PortTo   int
}

type Pod struct {
	Name      string
	Namespace string
}

// StopFn knows how to stop things
type StopFn func() error

type Poder interface {
	// PodExec executes a script within a specified pod
	PodExec(PodExecOpts, ...string) error
	// PortForward opens a port forwarding connection and returns a function to close that connection
	PortForward(PortForwardOpts) (StopFn, error)
	// PodReady returns a boolean indicating if the pod is ready or not
	PodReady(Pod) (bool, error)
}

type Client interface {
	Poder
	// Apply applies a manifest to the contextual cluster
	Apply(manifest io.Reader) error
}
