package certbot

import (
	"fmt"
)

const (
	apiVersion                  = "cert-manager.io/v1"
	clusterIssuerKind           = "ClusterIssuer"
	letsEncryptProductionServer = "https://acme-v02.api.letsencrypt.org/directory"
)

func newLetsEncryptClusterIssuer(email string) clusterIssuer {
	const name = "letsencrypt"

	return clusterIssuer{
		objectMeta: objectMeta{
			APIVersion: apiVersion,
			Kind:       clusterIssuerKind,
		},
		Metadata: clusterIssuerTypeMeta{
			Name:      name,
			Namespace: "kube-system",
		},
		Spec: clusterIssuerSpec{
			ACME: clusterIssuerSpecACME{
				Email:  email,
				Server: letsEncryptProductionServer,
				PrivateKeySecretRef: clusterIssuerSpecACMEPrivateKeySecretRef{
					Name: fmt.Sprintf("%s-key", name),
				},
				Solvers: []clusterIssuerSpecACMESolver{
					{
						HTTP01: clusterIssuerSpecACMESolverType{
							Ingress: clusterIssuerSpecACMESolverIngress{Class: "nginx"},
						},
					},
				},
			},
		},
	}
}

type objectMeta struct {
	APIVersion string `json:"apiVersion"`
	Kind       string `json:"kind"`
}

type clusterIssuerTypeMeta struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type clusterIssuerSpecACMEPrivateKeySecretRef struct {
	Name string `json:"name"`
}

type clusterIssuerSpecACMESolverIngress struct {
	Class string `json:"class"`
}

type clusterIssuerSpecACMESolverType struct {
	Ingress clusterIssuerSpecACMESolverIngress `json:"ingress"`
}

type clusterIssuerSpecACMESolver struct {
	HTTP01 clusterIssuerSpecACMESolverType `json:"http01"`
}

type clusterIssuerSpecACME struct {
	Email               string                                   `json:"email"`
	Server              string                                   `json:"server"`
	PrivateKeySecretRef clusterIssuerSpecACMEPrivateKeySecretRef `json:"privateKeySecretRef"`
	Solvers             []clusterIssuerSpecACMESolver            `json:"solvers"`
}

type clusterIssuerSpec struct {
	ACME clusterIssuerSpecACME `json:"acme"`
}

type clusterIssuer struct {
	objectMeta
	Metadata clusterIssuerTypeMeta `json:"metadata"`
	Spec     clusterIssuerSpec     `json:"spec"`
}
