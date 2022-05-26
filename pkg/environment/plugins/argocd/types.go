package argocd

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

const logFeature = "plugin/argocd"

type reconciler struct {
	cloudProvider cloud.Provider
}

type keyPair struct {
	PrivateKey []byte
	PublicKey  []byte
}

type repositorySecretOpts struct {
	SecretName           string
	RepositoryName       string
	RepositoryURI        string
	RepositoryPrivateKey string
	OperationsNamespace  string
}

type repository struct {
	URL string
}
