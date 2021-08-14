package linode

import (
	"context"
	"fmt"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/deifyed/xctl/pkg/config"
	"github.com/pkg/errors"
)

func (p *provider) HasPrimaryDomain(ctx context.Context, fqdn string) (bool, error) {
	domain := cloud.Domain{Host: fqdn}

	err := domain.Validate()
	if err != nil {
		return false, fmt.Errorf("validating domain: %w", err)
	}

	_, err = p.getLinodeDomain(ctx, domain.PrimaryDomain())
	if err != nil {
		if errors.Is(err, config.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("getting Linode domain: %w", err)
	}

	return true, nil
}
