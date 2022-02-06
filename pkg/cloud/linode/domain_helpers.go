package linode

import (
	"context"
	"fmt"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/linode/linodego"
)

func (p *provider) getLinodeDomain(ctx context.Context, domainName string) (linodego.Domain, error) {
	domains, err := p.client.ListDomains(ctx, &linodego.ListOptions{})
	if err != nil {
		return linodego.Domain{}, fmt.Errorf("listing domains: %w", err)
	}

	for _, domain := range domains {
		if domain.Domain == domainName {
			return domain, nil
		}
	}

	return linodego.Domain{}, config.ErrNotFound
}

func (p *provider) getLinodeDomainRecord(
	ctx context.Context,
	domainID int,
	subdomainName string,
) (linodego.DomainRecord, error) {
	if subdomainName == "" {
		subdomainName = "*"
	}

	records, err := p.client.ListDomainRecords(ctx, domainID, &linodego.ListOptions{})
	if err != nil {
		return linodego.DomainRecord{}, fmt.Errorf("listing domain records for ID %d: %w", domainID, err)
	}

	for _, record := range records {
		if record.Name == subdomainName {
			return record, nil
		}
	}

	return linodego.DomainRecord{}, config.ErrNotFound
}
