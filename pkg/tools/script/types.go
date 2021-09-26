package script

import "github.com/spf13/afero"

type Runner struct {
	fs  *afero.Afero
	env map[string]string
}
