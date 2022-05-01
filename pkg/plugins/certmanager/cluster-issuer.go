package certmanager

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

func newClusterIssuers(email string) (io.Reader, error) {
	t, err := template.New("issuer").Parse(issuerTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}
	opts := issuerOpts{Email: email}

	err = t.Execute(&buf, opts)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

//go:embed cluster-issuers.yaml
var issuerTemplate string

type issuerOpts struct {
	Email string `json:"email"`
}
