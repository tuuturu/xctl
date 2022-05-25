package application

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"text/template"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

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
