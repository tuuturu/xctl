package cloud

import (
	"context"
)

// Domain represents a domain with or without subdomains
type Domain struct {
	// Host defines the hostname of the domain
	Host string
}

type SubdomainServiceCRUDer interface {
	// CreateSubdomain knows how to create a subdomain in the cloud provider
	CreateSubdomain(ctx context.Context, domain Domain, target string) (Domain, error)
	// DeleteSubdomain knows how to delete a subdomain in the cloud provider
	DeleteSubdomain(ctx context.Context, domain Domain) error
	// GetSubdomain knows how to retrieve details about a domain in the cloud provider
	GetSubdomain(ctx context.Context, domain Domain) (Domain, error)
	// HasSubdomain knows if a subdomain exists in the cloud provider or not
	HasSubdomain(ctx context.Context, domain Domain) (bool, error)
}

type PrimaryDomainCRUDer interface {
	// HasPrimaryDomain knows if a domain exists in the cloud provider or not
	HasPrimaryDomain(ctx context.Context, domain Domain) (bool, error)
}

type DomainService interface {
	PrimaryDomainCRUDer
	SubdomainServiceCRUDer
}
