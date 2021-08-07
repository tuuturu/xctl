package handlers

import (
	"fmt"
	"io"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"

	"sigs.k8s.io/yaml"
)

func HandleApplication(_ io.Writer, _ bool, content []byte) error {
	var manifest v1alpha1.Application

	err := yaml.Unmarshal(content, &manifest)
	if err != nil {
		return fmt.Errorf("parsing application manifest: %w", err)
	}

	println(fmt.Sprintf("finished handling %s", manifest.Metadata.Name))

	return nil
}
