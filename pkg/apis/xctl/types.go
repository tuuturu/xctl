package xctl

import "io"

// IOStreams contains streams for interacting with the outside world
type IOStreams struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}
