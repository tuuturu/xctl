package cloud

type Provider interface {
	Authenticate() error
	ClusterService
}
