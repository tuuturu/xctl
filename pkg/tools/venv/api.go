package venv

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/afero"
)

const (
	defaultUsernameEnvKey         = "USER"
	defaultPasswdPath             = "/etc/passwd"
	defaultPasswdShellColumnIndex = 6
)

func GetCurrentShell(fs *afero.Afero) (string, error) {
	currentUser := os.Getenv(defaultUsernameEnvKey)

	passwdFile, err := fs.Open(defaultPasswdPath)
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", defaultPasswdPath, err)
	}

	defer func() {
		_ = passwdFile.Close()
	}()

	scanner := bufio.NewScanner(passwdFile)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, currentUser) {
			return strings.Split(line, ":")[defaultPasswdShellColumnIndex], nil
		}
	}

	return "", fmt.Errorf("finding line in passwd file: %w", err)
}

func MergeVariables(varSlices ...[]string) []string {
	result := map[string]string{}

	for _, slice := range varSlices {
		for key, value := range SliceAsMap(slice) {
			result[key] = value
		}
	}

	return MapAsSlice(result)
}

func MapAsSlice(m map[string]string) []string {
	result := make([]string, len(m))
	index := 0

	for key, value := range m {
		result[index] = fmt.Sprintf("%s=%s", key, value)

		index++
	}

	return result
}

func SliceAsMap(slice []string) map[string]string {
	m := make(map[string]string)

	for _, env := range slice {
		split := strings.SplitN(env, "=", 2)
		key := split[0]
		val := split[1]
		m[key] = val
	}

	return m
}
