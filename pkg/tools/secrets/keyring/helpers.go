package keyring

import (
	"fmt"

	"github.com/99designs/keyring"
	"github.com/deifyed/xctl/pkg/config"
)

func (c Client) open() (keyring.Keyring, error) {
	serviceName := fmt.Sprintf("%s-%s", serviceNamePrefix, c.EnvironmentName)

	environmentDirectory, err := config.GetAbsoluteXCTLClusterDir(c.EnvironmentName)
	if err != nil {
		return nil, fmt.Errorf("acquiring environment directory: %w", err)
	}

	return keyring.Open(keyring.Config{
		ServiceName:   serviceName,
		KeychainName:  serviceName,
		FileDir:       environmentDirectory,
		KWalletFolder: environmentDirectory,
		PassDir:       environmentDirectory,
		PassPrefix:    serviceName,
		WinCredPrefix: serviceName,
	})
}
