package helm

import "github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

// Client knows how to do operations with Helm on the cluster
type Client interface {
	// Install knows how to install a chart in the cluster
	Install(plugin v1alpha1.Plugin) error
	// Delete knows how to uninstall an existing release
	Delete(plugin v1alpha1.Plugin) error
	// Exists knows if a release exists in the cluster or not
	Exists(plugin v1alpha1.Plugin) (bool, error)
}
