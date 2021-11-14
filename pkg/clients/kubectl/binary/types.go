package binary

const (
	logFeature             = "kubectl/binary"
	kubeConfigPathKey      = "KUBECONFIG"
	portForwardWaitSeconds = 5
)

type kubectlBinaryClient struct {
	env         map[string]string
	kubectlPath string
}
