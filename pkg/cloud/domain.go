package cloud

import (
	"context"
)

type Domain struct {
	Host string
}

type SubdomainServiceCRUDer interface {
	CreateSubdomain(ctx context.Context, domain Domain, target string) (Domain, error)
	DeleteSubdomain(ctx context.Context, domain Domain) error
	GetSubdomain(ctx context.Context, domain Domain) (Domain, error)
	HasSubdomain(ctx context.Context, domain Domain) (bool, error)
}

type PrimaryDomainCRUDer interface {
	HasPrimaryDomain(ctx context.Context, domain Domain) (bool, error)
}

type DomainService interface {
	PrimaryDomainCRUDer
	SubdomainServiceCRUDer
}
