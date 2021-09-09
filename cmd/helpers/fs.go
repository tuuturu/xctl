package helpers

import (
	"fmt"

	"github.com/spf13/afero"
)

// CopyToFs copies a file from sourceFs to destinationFs
func CopyToFs(sourceFs *afero.Afero, destinationFs *afero.Afero, sourcePath, destinationPath string) error {
	content, err := sourceFs.ReadFile(sourcePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	err = destinationFs.WriteFile(destinationPath, content, 0o744)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}
