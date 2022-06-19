package kustomize

import (
	"fmt"
	"path"

	"github.com/spf13/afero"
)

// AddResourceToKustomization knows how to add a resource entry to a kustomization file
func AddResourceToKustomization(fs *afero.Afero, absoluteWorkDir string, resources ...string) error {
	kustomizationPath := path.Join(absoluteWorkDir, defaultKustomizationFilename)

	kustomizationFile, err := ensureReadKustomizationFile(fs, kustomizationPath)
	if err != nil {
		return fmt.Errorf("acquiring file: %w", err)
	}

	for _, resource := range resources {
		if contains(kustomizationFile.Resources, resource) {
			continue
		}

		kustomizationFile.Resources = append(kustomizationFile.Resources, resource)
	}

	err = ensureWriteKustomizationFile(fs, kustomizationPath, kustomizationFile)
	if err != nil {
		return fmt.Errorf("saving: %w", err)
	}

	return nil
}
