package binary

const kubeConfigPathKey = "KUBECONFIG"

type kubectlBinaryClient struct {
	kubectlPath string
	env         map[string]string
}
