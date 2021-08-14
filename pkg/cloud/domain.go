package cloud

import (
	"context"
)

type Domain struct {
	Host string
}

type SubdomainServiceCRUDer interface {
	CreateSubdomain(ctx context.Context, fqdn string) (Domain, error)
	DeleteSubdomain(ctx context.Context, fqdn string) error
	GetSubdomain(ctx context.Context, fqdn string) (Domain, error)
	HasSubdomain(ctx context.Context, fqdn string) (bool, error)
}

type PrimaryDomainCRUDer interface {
	HasPrimaryDomain(ctx context.Context, fqdn string) (bool, error)
}

type DomainService interface {
	PrimaryDomainCRUDer
	SubdomainServiceCRUDer
}
