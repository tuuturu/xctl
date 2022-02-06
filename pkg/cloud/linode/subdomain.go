package linode

import (
	"context"
	"fmt"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/linode/linodego"
	"github.com/pkg/errors"

	"github.com/deifyed/xctl/pkg/cloud"
)

func (p *provider) CreateSubdomain(ctx context.Context, domain cloud.Domain, target string) (cloud.Domain, error) {
	primaryDomain, err := p.getLinodeDomain(ctx, domain.PrimaryDomain())
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return cloud.Domain{}, fmt.Errorf("finding parent domain %s. Please register this domain with Linode: %w",
				domain.PrimaryDomain(),
				err,
			)
		}

		return cloud.Domain{}, fmt.Errorf("getting parent domain %s: %w", domain.PrimaryDomain(), err)
	}

	_, err = p.client.CreateDomainRecord(ctx, primaryDomain.ID, linodego.DomainRecordCreateOptions{
		Type:   linodego.RecordTypeA,
		Name:   domain.Subdomain(),
		Target: target,
		TTLSec: config.DefaultSubdomainTTLSeconds,
	})
	if err != nil {
		return cloud.Domain{}, fmt.Errorf("creating record: %w", err)
	}

	return domain, nil
}

func (p *provider) DeleteSubdomain(ctx context.Context, domain cloud.Domain) error {
	primaryDomain, err := p.getLinodeDomain(ctx, domain.PrimaryDomain())
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return fmt.Errorf("finding parent domain %s. Please register this domain with Linode: %w",
				domain.PrimaryDomain(),
				err,
			)
		}

		return fmt.Errorf("getting parent domain %s: %w", domain.PrimaryDomain(), err)
	}

	record, err := p.getLinodeDomainRecord(ctx, primaryDomain.ID, domain.Subdomain())
	if err != nil {
		return fmt.Errorf("getting domain record for %s: %w", domain.Subdomain(), err)
	}

	err = p.client.DeleteDomainRecord(ctx, primaryDomain.ID, record.ID)
	if err != nil {
		return fmt.Errorf("deleting domain name record: %w", err)
	}

	return nil
}

func (p *provider) GetSubdomain(ctx context.Context, domain cloud.Domain) (cloud.Domain, error) {
	linodeDomain, err := p.getLinodeDomain(ctx, domain.PrimaryDomain())
	if err != nil {
		return cloud.Domain{}, fmt.Errorf("getting primary domain: %w", err)
	}

	_, err = p.getLinodeDomainRecord(ctx, linodeDomain.ID, domain.Subdomain())
	if err != nil {
		return cloud.Domain{}, fmt.Errorf("getting record: %w", err)
	}

	return domain, nil
}

func (p *provider) HasSubdomain(ctx context.Context, domain cloud.Domain) (bool, error) {
	_, err := p.GetSubdomain(ctx, domain)
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("getting subdomain: %w", err)
	}

	return true, nil
}
