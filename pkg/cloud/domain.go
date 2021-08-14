package cloud

import (
	"context"
	"fmt"
	"strings"
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

type DomainService interface {
	SubdomainServiceCRUDer
}

// PrimaryDomain knows how to extract the second and top level domain i.e. tuuturu.org in dev.tuuturu.org
func (d Domain) PrimaryDomain() string {
	parts := strings.Split(d.Host, ".")

	topLevelDomain := parts[len(parts)-1]
	secondLevelDomain := parts[len(parts)-2]

	return fmt.Sprintf("%s.%s", secondLevelDomain, topLevelDomain)
}

// Subdomain knows how to extract the lowest level of domain i.e. dev in dev.tuuturu.org
func (d Domain) Subdomain() string {
	parts := strings.Split(d.Host, ".")

	return parts[0]
}

// FQDN ensures the Host provided satisfies the fully qualified domain name format
func (d Domain) FQDN() string {
	if strings.HasSuffix(d.Host, ".") {
		return d.Host
	}

	return fmt.Sprintf("%s.", d.Host)
}

// String returns the fully qualified domain name excluding the dot
func (d Domain) String() string {
	return strings.TrimRight(d.FQDN(), ".")
}
