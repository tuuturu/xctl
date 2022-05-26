package argocd

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"
	"text/template"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"github.com/deifyed/xctl/pkg/config"
)

//go:embed application-template.yaml
var applicationTemplate string

func scaffoldArgoCDApplication(environment v1alpha1.Environment, application v1alpha1.Application) (io.Reader, error) {
	t, err := template.New("application").Parse(applicationTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	buf := bytes.Buffer{}

	targetDirectory := path.Join(
		config.DefaultInfrastructureDir,
		config.DefaultApplicationsDir,
		application.Metadata.Name,
		config.DefaultApplicationsOverlaysDir,
		environment.Metadata.Name,
	)

	err = t.Execute(&buf, struct {
		ApplicationName      string
		ApplicationNamespace string
		OperationsNamespace  string
		TargetDirectory      string
		RepositoryURI        string
	}{
		ApplicationName:      application.Metadata.Name,
		ApplicationNamespace: application.Metadata.Namespace,
		OperationsNamespace:  config.DefaultOperationsNamespace,
		TargetDirectory:      targetDirectory,
		RepositoryURI:        environment.Spec.Repository,
	})
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}
