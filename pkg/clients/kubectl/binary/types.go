package binary

import "github.com/sirupsen/logrus"

const kubeConfigPathKey = "KUBECONFIG"

type kubectlBinaryClient struct {
	logger      *logrus.Logger
	env         map[string]string
	kubectlPath string
}
