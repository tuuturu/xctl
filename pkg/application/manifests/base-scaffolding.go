package manifests

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"path"
	"text/template"

	"sigs.k8s.io/yaml"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
)

func writeBaseManifests(fs readerWriter, targetDir string, application v1alpha1.Application) error {
	resources := make([]string, 0)

	deployment, err := scaffoldDeployment(application)
	if err != nil {
		return fmt.Errorf("scaffolding deployment: %w", err)
	}

	deploymentFilename := "deployment.yaml"
	resources = append(resources, deploymentFilename)

	err = fs.WriteReader(path.Join(targetDir, deploymentFilename), deployment)
	if err != nil {
		return fmt.Errorf("writing deployment: %w", err)
	}

	if requiresNetworking(application) {
		service, err := scaffoldService(application)
		if err != nil {
			return fmt.Errorf("scaffolding service: %w", err)
		}

		serviceFilename := "service.yaml"
		resources = append(resources, serviceFilename)

		err = fs.WriteReader(path.Join(targetDir, serviceFilename), service)
		if err != nil {
			return fmt.Errorf("writing service: %w", err)
		}
	}

	if requiresIngress(application) {
		ingress, err := scaffoldIngress(application)
		if err != nil {
			return fmt.Errorf("scaffolding ingress: %w", err)
		}

		ingressFilename := "ingress.yaml"
		resources = append(resources, ingressFilename)

		err = fs.WriteReader(path.Join(targetDir, ingressFilename), ingress)
		if err != nil {
			return fmt.Errorf("writing ingress: %w", err)
		}
	}

	rawKustomization, err := yaml.Marshal(&kustomize{Resources: resources})
	if err != nil {
		return fmt.Errorf("marshalling kustomization file: %w", err)
	}

	err = fs.WriteReader(path.Join(targetDir, "kustomization.yaml"), bytes.NewReader(rawKustomization))
	if err != nil {
		return fmt.Errorf("writing kustomization: %w", err)
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

func requiresNetworking(application v1alpha1.Application) bool {
	return application.Spec.Port != ""
}

func requiresIngress(application v1alpha1.Application) bool {
	return requiresNetworking(application) && application.Spec.Url != ""
}
