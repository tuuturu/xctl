package cloud

// Provider defines required functionality xctl expects from a cloud provider
type Provider interface {
	// Authenticate knows how to authenticate to a cloud provider
	Authenticate() error
	ClusterService
	DomainService
}
