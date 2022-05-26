package manifests

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"
	"text/template"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func writeBaseManifests(fs readerWriter, targetDir string, application v1alpha1.Application) error {
	deployment, err := scaffoldDeployment(application)
	if err != nil {
		return fmt.Errorf("scaffolding deployment: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "deployment.yaml"), deployment)
	if err != nil {
		return fmt.Errorf("writing deployment: %w", err)
	}

	service, err := scaffoldService(application)
	if err != nil {
		return fmt.Errorf("scaffolding service: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "service.yaml"), service)
	if err != nil {
		return fmt.Errorf("writing service: %w", err)
	}

	ingress, err := scaffoldIngress(application)
	if err != nil {
		return fmt.Errorf("scaffolding ingress: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "ingress.yaml"), ingress)
	if err != nil {
		return fmt.Errorf("writing ingress: %w", err)
	}

	return nil
}

//go:embed templates/deployment.yaml
var deploymentTemplate string

func scaffoldDeployment(application v1alpha1.Application) (io.Reader, error) {
	t, err := template.New("deployment").Parse(deploymentTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		ApplicationName string
		ImageURI        string
	}{
		ApplicationName: application.Metadata.Name,
		ImageURI:        application.Spec.Image,
	})
	if err != nil {
		return nil, fmt.Errorf("executing: %w", err)
	}

	return &buf, nil
}

//go:embed templates/service.yaml
var serviceTemplate string

func scaffoldService(application v1alpha1.Application) (io.Reader, error) {
	t, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		ApplicationName string
		ApplicationPort string
	}{
		ApplicationName: application.Metadata.Name,
		ApplicationPort: application.Spec.Port,
	})
	if err != nil {
		return nil, fmt.Errorf("executing: %w", err)
	}

	return &buf, nil
}

//go:embed templates/ingress.yaml
var ingressTemplate string

func scaffoldIngress(application v1alpha1.Application) (io.Reader, error) {
	t, err := template.New("ingress").Parse(ingressTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		ApplicationName string
		Host            string
	}{
		ApplicationName: application.Metadata.Name,
		Host:            application.Spec.Url,
	})
	if err != nil {
		return nil, fmt.Errorf("executing: %w", err)
	}

	return &buf, nil
}
