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

// Poder defines operations done on pods
type Poder interface {
	// PodExec executes a script within a specified pod
	PodExec(PodExecOpts, ...string) error
	// PortForward opens a port forwarding connection and returns a function to close that connection
	PortForward(PortForwardOpts) (StopFn, error)
	// PodReady returns a boolean indicating if the pod is ready or not
	PodReady(Pod) (bool, error)
}

type Resourcer interface {
	// Apply applies a manifest to the contextual cluster
	Apply(manifest io.Reader) error
	// Delete removes a manifest from teh contextual cluster
	Delete(manifest io.Reader) error
	// Get retrieves a named resource of a certain type from a specific namespace
	Get(namespace string, resourceType string, name string) (io.Reader, error)
	// DeleteResource removes a named resource of a certain kind from a specific namespace
	DeleteResource(namespace string, kind string, name string) error
	// IsReady knows if a Kubernetes resource is ready or not
	IsReady(Selector) (bool, error)
}

type Operator interface {
	// GetUserToken retrieves the authenticated user's token
	GetUserToken() (io.Reader, error)
}

// Client defines operations done on a Kubernetes cluster
type Client interface {
	Poder
	Resourcer
	Operator
}

// Selector describes a resource in the cluster
type Selector struct {
	// Namespace defines what namespace the resource resides in
	Namespace string
	// Kind defines what type of resource it is
	Kind string
	// Name defines the name of the resource
	Name string
}

// PodExecOpts defines required data for executing commands on a pod
type PodExecOpts struct {
	Pod    Pod
	Stdout io.Writer
}

// PortForwardOpts defines required data for forwarding a port from a service
type PortForwardOpts struct {
	Service Service

	ServicePort int
	LocalPort   int
}

// Pod defines required data for identifying a pod
type Pod struct {
	Name      string
	Namespace string
}

// Service defines required data for identifying a service
type Service struct {
	Name      string
	Namespace string
}

// StopFn knows how to stop things
type StopFn func() error
