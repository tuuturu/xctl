package binary

import (
	"github.com/deifyed/xctl/pkg/tools/clients/kubectl"
)

func getCurrentContext(cfg kubeConfig, name string) (kubeConfigContextContext, error) {
	for _, ctx := range cfg.Contexts {
		if ctx.Name == name {
			return ctx.Context, nil
		}
	}

	return kubeConfigContextContext{}, kubectl.ErrNotFound
}

func getUserForContext(cfg kubeConfig, ctx kubeConfigContextContext) (kubeConfigUsersUser, error) {
	for _, user := range cfg.Users {
		if user.Name == ctx.User {
			return user.User, nil
		}
	}

	return kubeConfigUsersUser{}, kubectl.ErrNotFound
}
