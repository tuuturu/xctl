package binary

import "fmt"

func (k kubectlBinaryClient) envAsArray() []string {
	env := make([]string, len(k.env))
	index := 0

	for key, value := range k.env {
		env[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return env
}
