package linode

import (
	"context"
	"fmt"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/pkg/errors"
)

func (p *provider) HasPrimaryDomain(ctx context.Context, domain cloud.Domain) (bool, error) {
	err := domain.Validate()
	if err != nil {
		return false, fmt.Errorf("validating domain: %w", err)
	}

	_, err = p.getLinodeDomain(ctx, domain.PrimaryDomain())
	if err != nil {
		if errors.Is(err, cloud.ErrNotFound) {
			return false, nil
		}

		return false, fmt.Errorf("getting Linode domain: %w", err)
	}

	return true, nil
}
