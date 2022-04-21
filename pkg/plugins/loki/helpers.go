package loki

import (
	"fmt"
	"io"
	"strings"
)

func handleManifests(fn func(reader io.Reader) error, manifests []string) error {
	for _, manifest := range manifests {
		err := fn(strings.NewReader(manifest))
		if err != nil {
			return fmt.Errorf("handling manifest: %w", err)
		}
	}

	return nil
}
