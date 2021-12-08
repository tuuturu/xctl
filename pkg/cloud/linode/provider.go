package linode

import (
	"context"
	"fmt"
	"time"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/deifyed/xctl/pkg/config"
	"github.com/linode/linodego"
	"github.com/pkg/errors"
)

func (p *provider) getCluster(ctx context.Context, clusterName string) (linodego.LKECluster, error) {
	clusters, err := p.client.ListLKEClusters(ctx, &linodego.ListOptions{})
	if err != nil {
		return linodego.LKECluster{}, fmt.Errorf("retrieving existing LKE clusters: %w", err)
	}

	for _, cluster := range clusters {
		if cluster.Label == clusterName {
			return cluster, nil
		}
	}

	return linodego.LKECluster{}, config.ErrNotFound
}

func (p *provider) await(test pollTestFn) (err error) {
	timeout := time.Now().Add(defaultTimeoutMinutes * time.Minute)
	delayFunction := func() { time.Sleep(defaultDelaySeconds * time.Second) }

	var ready bool

	for !ready {
		delayFunction()

		if time.Now().After(timeout) {
			return config.ErrTimeout
		}

		ready, err = test()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *provider) awaitCreation(ctx context.Context, clusterID int) error {
	log := logging.GetLogger(logFeature, "awaitCreation")

	return p.await(func() (bool, error) {
		ok, err := nodePoolCheck(ctx, p.client, clusterID)
		if err != nil {
			return false, fmt.Errorf("checking node pools: %w", err)
		}

		if !ok {
			return false, nil
		}

		log.Debug("node pools are ready")

		return true, nil
	})
}

func (p *provider) awaitDeletion(ctx context.Context, clusterName string) error {
	return p.await(func() (bool, error) {
		_, err := p.getCluster(ctx, clusterName)
		if err != nil {
			if errors.Is(err, config.ErrNotFound) {
				return true, nil
			}

			return false, fmt.Errorf("getting LKE cluster: %w", err)
		}

		return false, nil
	})
}

func nodePoolCheck(ctx context.Context, client linodego.Client, clusterID int) (bool, error) {
	pools, err := client.ListLKEClusterPools(ctx, clusterID, &linodego.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("listing LKE cluster pools: %w", err)
	}

	for _, pool := range pools {
		for _, node := range pool.Linodes {
			if node.Status != linodego.LKELinodeReady {
				return false, nil
			}
		}
	}

	return true, nil
}
