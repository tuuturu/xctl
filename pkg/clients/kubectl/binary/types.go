package binary

const (
	logFeature             = "kubectl/binary"
	kubeConfigPathKey      = "KUBECONFIG"
	portforwardWaitSeconds = 2
)

type kubectlBinaryClient struct {
	env         map[string]string
	kubectlPath string
}
