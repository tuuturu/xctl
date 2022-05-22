package keyring

import "fmt"

func generateServiceName(environmentName string) string {
	return fmt.Sprintf("%s-%s", serviceNamePrefix, environmentName)
}
