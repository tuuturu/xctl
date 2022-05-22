package linode

import (
	"github.com/deifyed/xctl/pkg/cloud"
)

func NewLinodeProvider() cloud.Provider {
	return &provider{}
}
