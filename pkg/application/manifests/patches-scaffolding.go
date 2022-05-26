package manifests

import (
	"bytes"
	"fmt"
	"path"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

func writeOverlaysPatches(fs *afero.Afero, targetDir string, _ v1alpha1.Application) error {
	kustomizationFile := kustomize{Resources: []string{"../../base"}}

	rawKustomizationFile, err := yaml.Marshal(&kustomizationFile)
	if err != nil {
		return fmt.Errorf("marshalling kustomization file: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, kustomizationFilename), bytes.NewReader(rawKustomizationFile))
	if err != nil {
		return fmt.Errorf("writing kustomization file: %w", err)
	}

	return nil
}
