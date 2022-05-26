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

func establishConfiguration(fs *afero.Afero, absoluteRepositoryRootDirectory string, environment v1alpha1.Environment) (io.Reader, error) {
	absoluteArgoCDConfigDir := configDir(absoluteRepositoryRootDirectory, environment.Metadata.Name)
	absoluteApplicationsDir := applicationsDir(absoluteRepositoryRootDirectory, environment.Metadata.Name)

	err := fs.MkdirAll(absoluteApplicationsDir, 0o700)
	if err != nil {
		return nil, fmt.Errorf("preparing directories: %w", err)
	}

	appManifest, err := buildApplicationsApplication(buildApplicationsApplicationOpts{
		OperationsNamespace: config.DefaultOperationsNamespace,
		ApplicationsDir:     applicationsDir("", environment.Metadata.Name),
		RepositoryURI:       environment.Spec.Repository,
	})
	if err != nil {
		return nil, fmt.Errorf("building applications application: %w", err)
	}

	rawAppManifest, err := io.ReadAll(appManifest)
	if err != nil {
		return nil, fmt.Errorf("buffering app manifest: %w", err)
	}

	err = fs.WriteReader(path.Join(absoluteArgoCDConfigDir, applicationsApplicationsFilename), bytes.NewReader(rawAppManifest))
	if err != nil {
		return nil, fmt.Errorf("writing applications application: %w", err)
	}

	applicationsReadme, err := buildApplicationsReadme()
	if err != nil {
		return nil, fmt.Errorf("building applications readme: %w", err)
	}

	err = fs.WriteReader(path.Join(absoluteApplicationsDir, readmeFilename), applicationsReadme)
	if err != nil {
		return nil, fmt.Errorf("writing applications readme: %w", err)
	}

	return bytes.NewReader(rawAppManifest), nil
}

func buildApplicationsApplication(opts buildApplicationsApplicationOpts) (io.Reader, error) {
	buf := bytes.Buffer{}

	t, err := template.New("applicationsApplication").Parse(applicationsApplicationTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, opts)
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

func buildApplicationsReadme() (io.Reader, error) {
	buf := bytes.Buffer{}

	t, err := template.New("applicationsApplication").Parse(applicationsReadmeTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing template: %w", err)
	}

	err = t.Execute(&buf, struct{}{})
	if err != nil {
		return nil, fmt.Errorf("executing template: %w", err)
	}

	return &buf, nil
}

func configDir(root string, environmentName string) string {
	return path.Join(root, config.DefaultInfrastructureDir, environmentName, pluginName)
}

func applicationsDir(root string, environmentName string) string {
	return path.Join(configDir(root, environmentName), "applications")
}

const (
	applicationsApplicationsFilename = "applications.yaml"
	readmeFilename                   = "README.md"
)

type buildApplicationsApplicationOpts struct {
	OperationsNamespace string
	ApplicationsDir     string
	RepositoryURI       string
}

var (
	//go:embed templates/applications-application.yaml
	applicationsApplicationTemplate string
	//go:embed templates/applications-readme.md
	applicationsReadmeTemplate string
)
