package binary

const (
	logFeature             = "kubectl/binary"
	kubeConfigPathKey      = "KUBECONFIG"
	portforwardWaitSeconds = 1
)

type kubectlBinaryClient struct {
	env         map[string]string
	kubectlPath string
}
