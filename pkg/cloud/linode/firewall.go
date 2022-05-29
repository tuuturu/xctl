package linode

import (
	"context"
	"fmt"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/linode/linodego"
)

// Docs: https://www.linode.com/community/questions/19155/securing-k8s-cluster

func (p *provider) createClusterNodesFirewall(ctx context.Context, manifest v1alpha1.Environment, clusterID int) (int, error) {
	nodes, err := p.getClusterNodes(ctx, clusterID)
	if err != nil {
		return 0, fmt.Errorf("acquiring nodes: %w", err)
	}

	nodeIDs := make([]int, len(nodes))

	for index, node := range nodes {
		nodeIDs[index] = node.ID
	}

	firewall, err := p.client.CreateFirewall(ctx, linodego.FirewallCreateOptions{
		Label: clusterNodesFirewallName(manifest),
		Rules: linodego.FirewallRuleSet{
			InboundPolicy: policyDrop,
			Inbound: []linodego.FirewallRule{
				{
					Action:      policyAccept,
					Label:       componentNamer(manifest, "fwr", "health"),
					Description: "Kubelet health checks",
					Ports:       "10250",
					Protocol:    linodego.TCP,
					Addresses:   linodego.NetworkAddresses{IPv4: &[]string{privateNetworkCIDR}},
				},
				{
					Action:      policyAccept,
					Label:       componentNamer(manifest, "fwr", "proxy"),
					Description: "Wireguard tunneling for kubectl proxy",
					Ports:       "51820",
					Protocol:    linodego.UDP,
					Addresses:   linodego.NetworkAddresses{IPv4: &[]string{privateNetworkCIDR}},
				},
				{
					Action:      policyAccept,
					Label:       componentNamer(manifest, "fwr", "calico"),
					Description: "Calico BGP traffic",
					Ports:       "179",
					Protocol:    linodego.TCP,
					Addresses:   linodego.NetworkAddresses{IPv4: &[]string{privateNetworkCIDR}},
				},
				{
					Action:      policyAccept,
					Label:       componentNamer(manifest, "fwr", "nb"),
					Description: "Allows traffic from nodebalancers",
					Ports:       "30000-32768",
					Protocol:    linodego.TCP,
					Addresses:   linodego.NetworkAddresses{IPv4: &[]string{privateNetworkCIDR}},
				},
			},
			OutboundPolicy: policyAccept,
		},
		Tags:    defaultLabels(manifest),
		Devices: linodego.DevicesCreationOptions{Linodes: nodeIDs},
	})
	if err != nil {
		return 0, fmt.Errorf("creating: %w", err)
	}

	return firewall.ID, nil
}

func (p *provider) createNodebalancerFirewall(ctx context.Context, manifest v1alpha1.Environment, clusterID int) (int, error) {
	nodes, err := p.getClusterNodes(ctx, clusterID)
	if err != nil {
		return 0, fmt.Errorf("acquiring nodes: %w", err)
	}

	nodebalancer, err := p.getClusterNodebalancer(ctx, nodes)
	if err != nil {
		return 0, fmt.Errorf("acquiring nodebalancer: %w", err)
	}

	firewall, err := p.client.CreateFirewall(ctx, linodego.FirewallCreateOptions{
		Label: nodeBalancerFirewallName(manifest),
		Rules: linodego.FirewallRuleSet{
			Inbound: []linodego.FirewallRule{
				{
					Action:      policyAccept,
					Label:       componentNamer(manifest, "fwr", "std"),
					Description: "Allow incoming requests from the internet on HTTPs",
					Ports:       "80, 443",
					Protocol:    linodego.TCP,
					Addresses: linodego.NetworkAddresses{
						IPv4: &[]string{"0.0.0.0/0"},
						IPv6: &[]string{"::/0"},
					},
				},
			},
			InboundPolicy:  policyDrop,
			OutboundPolicy: policyAccept,
		},
		Tags:    defaultLabels(manifest),
		Devices: linodego.DevicesCreationOptions{NodeBalancers: []int{nodebalancer.ID}},
	})
	if err != nil {
		return 0, fmt.Errorf("creating: %w", err)
	}

	return firewall.ID, nil
}

func (p *provider) deleteFirewall(ctx context.Context, name string) error {
	firewalls, err := p.client.ListFirewalls(ctx, &linodego.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing firewalls: %w", err)
	}

	firewallID := -1

	for _, firewall := range firewalls {
		if firewall.Label == name {
			firewallID = firewall.ID
		}
	}

	if firewallID == -1 {
		return nil
	}

	err = p.client.DeleteFirewall(ctx, firewallID)
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

func (p *provider) deleteClusterNodesFirewall(ctx context.Context, manifest v1alpha1.Environment) error {
	err := p.deleteFirewall(ctx, clusterNodesFirewallName(manifest))
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

func (p *provider) deleteNodebalancerFirewall(ctx context.Context, manifest v1alpha1.Environment) error {
	err := p.deleteFirewall(ctx, nodeBalancerFirewallName(manifest))
	if err != nil {
		return fmt.Errorf("deleting: %w", err)
	}

	return nil
}

func clusterNodesFirewallName(manifest v1alpha1.Environment) string {
	return componentNamer(manifest, "fw", "nodes")
}

func nodeBalancerFirewallName(manifest v1alpha1.Environment) string {
	return componentNamer(manifest, "fw", "nb")
}

const (
	policyAccept       = "ACCEPT"
	policyDrop         = "DROP"
	privateNetworkCIDR = "192.168.128.0/17"
)
