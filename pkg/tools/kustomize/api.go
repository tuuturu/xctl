package kustomize

import (
	"fmt"
	"path"

	"github.com/spf13/afero"
)

// AddResourceToKustomization knows how to add a resource entry to a kustomization file
func AddResourceToKustomization(fs *afero.Afero, absoluteWorkDir string, resourcePath string) error {
	kustomizationPath := path.Join(absoluteWorkDir, defaultKustomizationFilename)

	kustomizationFile, err := ensureReadKustomizationFile(fs, kustomizationPath)
	if err != nil {
		return fmt.Errorf("acquiring file: %w", err)
	}

	if contains(kustomizationFile.Resources, resourcePath) {
		return nil
	}

	kustomizationFile.Resources = append(kustomizationFile.Resources, resourcePath)

	err = ensureWriteKustomizationFile(fs, kustomizationPath, kustomizationFile)
	if err != nil {
		return fmt.Errorf("saving: %w", err)
	}

	return nil
}
