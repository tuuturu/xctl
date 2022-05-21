package keyring

import "fmt"

func generateServiceName(clusterName string) string {
	return fmt.Sprintf("%s-%s", serviceNamePrefix, clusterName)
}
