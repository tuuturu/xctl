package manifests

import "io"

type reconciler struct {
	absoluteApplicationDir string
}

// readerWriter simplifies requirements for Afero's fs.WriteReader func
type readerWriter interface {
	WriteReader(string, io.Reader) error
}
