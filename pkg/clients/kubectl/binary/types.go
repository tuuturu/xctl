package binary

const kubeConfigPathKey = "KUBECONFIG"

type kubectlBinaryClient struct {
	env         map[string]string
	kubectlPath string
}
