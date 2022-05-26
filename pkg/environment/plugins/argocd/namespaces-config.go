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
	"github.com/spf13/afero"
)

func establishNamespacesConfiguration(fs *afero.Afero, absoluteRepositoryRootDirectory string, environment v1alpha1.Environment) (io.Reader, error) {
	absoluteArgoCDConfigDir := configDir(absoluteRepositoryRootDirectory, environment.Metadata.Name)
	absoluteNamespacesDir := namespacesDir(absoluteRepositoryRootDirectory, environment.Metadata.Name)

	err := fs.MkdirAll(absoluteNamespacesDir, 0o700)
	if err != nil {
		return nil, fmt.Errorf("preparing directories: %w", err)
	}

	appManifest, err := buildNamespacesApplication(buildArgoCDApplicationOpts{
		OperationsNamespace: config.DefaultOperationsNamespace,
		TargetDirectory:     namespacesDir("", environment.Metadata.Name),
		RepositoryURI:       environment.Spec.Repository,
	})
	if err != nil {
		return nil, fmt.Errorf("building applications application: %w", err)
	}

	rawAppManifest, err := io.ReadAll(appManifest)
	if err != nil {
		return nil, fmt.Errorf("buffering app manifest: %w", err)
	}

	err = fs.WriteReader(path.Join(absoluteArgoCDConfigDir, namespacesApplicationFilename), bytes.NewReader(rawAppManifest))
	if err != nil {
		return nil, fmt.Errorf("writing applications application: %w", err)
	}

	namespacesReadme, err := buildNamespacesReadme()
	if err != nil {
		return nil, fmt.Errorf("building applications readme: %w", err)
	}

	err = fs.WriteReader(path.Join(absoluteNamespacesDir, readmeFilename), namespacesReadme)
	if err != nil {
		return nil, fmt.Errorf("writing applications readme: %w", err)
	}

	return bytes.NewReader(rawAppManifest), nil
}

func buildNamespacesApplication(opts buildArgoCDApplicationOpts) (io.Reader, error) {
	buf := bytes.Buffer{}

	t, err := template.New("namespaces-application").Parse(namespacesApplicationTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, opts)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

func buildNamespacesReadme() (io.Reader, error) {
	buf := bytes.Buffer{}

	t, err := template.New("readme").Parse(namespacesReadmeTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

func namespacesDir(root string, environmentName string) string {
	return path.Join(configDir(root, environmentName), "namespaces")
}

const namespacesApplicationFilename = "namespaces.yaml"

type buildArgoCDApplicationOpts struct {
	OperationsNamespace string
	TargetDirectory     string
	RepositoryURI       string
}

var (
	//go:embed templates/namespaces-application.yaml
	namespacesApplicationTemplate string
	//go:embed templates/namespaces-readme.md
	namespacesReadmeTemplate string
)
