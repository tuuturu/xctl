package script

import "fmt"

func mapAsSlice(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, value := range m {
		result[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return result
}
