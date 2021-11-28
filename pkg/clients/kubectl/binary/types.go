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

type commandLogFields struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}
