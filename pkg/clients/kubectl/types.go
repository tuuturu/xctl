package kubectl

import "io"

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

type Client interface {
	// PodExec executes a script within a specified pod
	PodExec(PodExecOpts, ...string) error
	// PortForward opens a port forwarding connection and returns a function to close that connection
	PortForward(PortForwardOpts) (StopFn, error)
}
