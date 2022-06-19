package kustomize

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/deifyed/xctl/pkg/tools/paths"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

func ensureReadKustomizationFile(fs *afero.Afero, kustomizationPath string) (file, error) {
	var kustomizationFile file

	rawKustomization, err := fs.ReadFile(kustomizationPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return file{}, nil
		}

		return file{}, fmt.Errorf("reading file: %w", err)
	}

	err = yaml.Unmarshal(rawKustomization, &kustomizationFile)
	if err != nil {
		return file{}, fmt.Errorf("unmarshalling: %w", err)
	}

	return kustomizationFile, nil
}

func ensureWriteKustomizationFile(fs *afero.Afero, kustomizationPath string, kustomizationFile file) error {
	err := fs.MkdirAll(path.Dir(kustomizationPath), paths.DefaultDirectoryPermissions)
	if err != nil {
		return fmt.Errorf("preparing work directory: %w", err)
	}

	payload, err := yaml.Marshal(kustomizationFile)
	if err != nil {
		return fmt.Errorf("marshalling: %w", err)
	}

	err = fs.WriteReader(kustomizationPath, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
