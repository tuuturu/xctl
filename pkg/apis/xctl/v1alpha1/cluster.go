package v1alpha1

import (
	"strings"

	"github.com/deifyed/xctl/pkg/config"
)

const ClusterKind = "Cluster"

type ClusterSpecPlugins struct {
	CertBot                bool `json:"certBot"`
	NginxIngressController bool `json:"nginxIngressController"`
	Vault                  bool `json:"vault"`
}

type ClusterSpec struct {
	RootDomain string             `json:"rootDomain"`
	AdminEmail string             `json:"adminEmail"`
	Plugins    ClusterSpecPlugins `json:"plugins"`
}

type Cluster struct {
	TypeMeta `json:",inline"`
	Metadata Metadata    `json:"metadata"`
	Spec     ClusterSpec `json:"spec"`
}

func NewDefaultCluster() Cluster {
	return Cluster{
		TypeMeta: TypeMeta{
			Kind:       ClusterKind,
			APIVersion: apiVersion,
		},
		Spec: ClusterSpec{
			Plugins: ClusterSpecPlugins{
				CertBot:                true,
				NginxIngressController: true,
				Vault:                  true,
			},
		},
	}
}

func (c Cluster) ComponentName(componentType string, id string) string {
	componentName := strings.Join([]string{config.ApplicationName, c.Metadata.Name, componentType, id}, "-")
	componentName = strings.ToLower(componentName)
	componentName = strings.TrimSuffix(componentName, "-")

	return componentName
}
