package manifests

import (
	"fmt"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/tools/kustomize"
	"github.com/spf13/afero"
)

func writeOverlaysPatches(fs *afero.Afero, targetDir string, _ v1alpha1.Application) error {
	err := kustomize.AddResourceToKustomization(fs, targetDir, "../../base")
	if err != nil {
		return fmt.Errorf("adding resource to kustomization file: %w", err)
	}

	return nil
}
