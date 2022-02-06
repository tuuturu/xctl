package linode

import (
	"context"
	"fmt"
	"strings"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/linode/linodego"
)

func (p *provider) getClusterNodebalancer(ctx context.Context, nodes []cloud.ClusterNode) (linodego.NodeBalancer, error) { //nolint:lll
	nodebalancers, err := p.client.ListNodeBalancers(ctx, &linodego.ListOptions{})
	if err != nil {
		return linodego.NodeBalancer{}, fmt.Errorf("listing node balancers: %w", err)
	}

	for _, nodebalancer := range nodebalancers {
		nodebalancerNodes, err := p.getNodebalancerTargetNodes(ctx, nodebalancer)
		if err != nil {
			return linodego.NodeBalancer{}, fmt.Errorf("getting target nodes for node balancer: %w", err)
		}

		if containsNodes(nodebalancerNodes, nodes) {
			return nodebalancer, nil
		}
	}

	return linodego.NodeBalancer{}, cloud.ErrNotFound
}

func (p *provider) getNodebalancerTargetNodes(ctx context.Context, nodebalancer linodego.NodeBalancer) ([]linodego.NodeBalancerNode, error) { //nolint:lll
	nodes := make(map[string]linodego.NodeBalancerNode)

	configs, err := p.client.ListNodeBalancerConfigs(ctx, nodebalancer.ID, &linodego.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing node balancer configs: %w", err)
	}

	for _, config := range configs {
		nodebalancerNodes, err := p.client.ListNodeBalancerNodes(ctx, nodebalancer.ID, config.ID, &linodego.ListOptions{})
		if err != nil {
			return nil, fmt.Errorf("listing node balancer nodes: %w", err)
		}

		for _, nodeBalancerNode := range nodebalancerNodes {
			nodes[strings.Split(nodeBalancerNode.Address, ":")[0]] = nodeBalancerNode
		}
	}

	nodesArray := make([]linodego.NodeBalancerNode, len(nodes))
	index := 0

	for _, value := range nodes {
		nodesArray[index] = value

		index++
	}

	return nodesArray, nil
}

func containsNodes(nodebalancerNodes []linodego.NodeBalancerNode, nodes []cloud.ClusterNode) bool {
	addresses := make([]string, len(nodebalancerNodes))

	for index, nodebalancerNode := range nodebalancerNodes {
		addresses[index] = strings.Split(nodebalancerNode.Address, ":")[0]
	}

	for _, node := range nodes {
		if containsNode(addresses, node.IPv4) {
			return true
		}
	}

	return false
}

func containsNode(haystack []string, needle string) bool {
	for _, potentialHit := range haystack {
		if potentialHit == needle {
			return true
		}
	}

	return false
}
