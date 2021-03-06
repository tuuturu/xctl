package linode

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"github.com/deifyed/xctl/pkg/cloud"

	"github.com/deifyed/xctl/pkg/tools/logging"

	"github.com/linode/linodego"
	"github.com/pkg/errors"
)

func (p *provider) getCluster(ctx context.Context, manifest v1alpha1.Environment) (linodego.LKECluster, error) {
	clusters, err := p.client.ListLKEClusters(ctx, &linodego.ListOptions{})
	if err != nil {
		return linodego.LKECluster{}, fmt.Errorf("retrieving existing LKE clusters: %w", err)
	}

	clusterName := componentNamer(manifest, componentTypeCluster, "")

	for _, cluster := range clusters {
		if cluster.Label == clusterName {
			return cluster, nil
		}
	}

	return linodego.LKECluster{}, cloud.ErrNotFound
}

func (p *provider) getClusterNodes(ctx context.Context, clusterID int) ([]cloud.ClusterNode, error) {
	pools, err := p.client.ListLKEClusterPools(ctx, clusterID, &linodego.ListOptions{})
	if err != nil {
		return []cloud.ClusterNode{}, err
	}

	nodes := make([]cloud.ClusterNode, 0)
	instances := make([]*linodego.Instance, 0)

	for _, pool := range pools {
		newInstances, err := p.getInstancesFromPool(ctx, pool)
		if err != nil {
			return nil, fmt.Errorf("retrieving instances from pool: %w", err)
		}

		instances = append(instances, newInstances...)
	}

	for _, instance := range instances {
		localIP, err := getLocalIP(instance.IPv4)
		if err != nil {
			return nil, fmt.Errorf("acquiring local IP for instance: %w", err)
		}

		nodes = append(nodes, cloud.ClusterNode{
			ID:   instance.ID,
			Name: instance.Label,
			IPv4: localIP,
		})
	}

	return nodes, nil
}

func (p *provider) getInstancesFromPool(ctx context.Context, pool linodego.LKEClusterPool) ([]*linodego.Instance, error) { //nolint:lll
	instances := make([]*linodego.Instance, len(pool.Linodes))

	for index, node := range pool.Linodes {
		instance, err := p.client.GetInstance(ctx, node.InstanceID)
		if err != nil {
			return nil, fmt.Errorf("getting instance: %w", err)
		}

		instances[index] = instance
	}

	return instances, nil
}

func (p *provider) await(test pollTestFn) (err error) {
	timeout := time.Now().Add(defaultTimeoutMinutes * time.Minute)
	delayFunction := func() { time.Sleep(defaultDelaySeconds * time.Second) }

	var ready bool

	for !ready {
		delayFunction()

		if time.Now().After(timeout) {
			return cloud.ErrTimeout
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

func (p *provider) awaitDeletion(ctx context.Context, manifest v1alpha1.Environment) error {
	return p.await(func() (bool, error) {
		_, err := p.getCluster(ctx, manifest)
		if err != nil {
			if errors.Is(err, cloud.ErrNotFound) {
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

func getLocalIP(ips []*net.IP) (string, error) {
	for _, ip := range ips {
		if strings.HasPrefix(ip.String(), "192.168.") {
			return ip.String(), nil
		}
	}

	return "", cloud.ErrNotFound
}
