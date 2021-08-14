package linode

import (
	"fmt"
	"net/http"
	"os"

	"github.com/deifyed/xctl/pkg/cloud"
	"github.com/linode/linodego"
	"golang.org/x/oauth2"
)

func (p *provider) Authenticate() error {
	apiKey, ok := os.LookupEnv(linodego.APIEnvVar)
	if !ok {
		return fmt.Errorf(fmt.Sprintf("finding Linode API token (%s)", linodego.APIEnvVar))
	}

	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: apiKey})

	oauth2Client := &http.Client{
		Transport: &oauth2.Transport{Source: tokenSource},
	}

	p.client = linodego.NewClient(oauth2Client)
	p.client.SetDebug(false)

	return nil
}

func NewLinodeProvider() cloud.Provider {
	return &provider{}
}
