package namespace

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"text/template"
)

//go:embed namespace-template.yaml
var namespaceTemplate string

func scaffoldNamespace(name string) (io.Reader, error) {
	t, err := template.New("namespace").Parse(namespaceTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		NamespaceName string
	}{NamespaceName: name})
	if err != nil {
		return nil, fmt.Errorf("executing: %w", err)
	}

	return &buf, nil
}
