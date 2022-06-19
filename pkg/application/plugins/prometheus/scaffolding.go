package prometheus

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"text/template"

	"github.com/deifyed/xctl/pkg/config"
)

//go:embed service-monitor-template.yaml
var serviceMonitorTemplate string

func scaffoldServiceMonitor(applicationName string, metricsPath string) (io.Reader, error) {
	t, err := template.New("service-monitor").Parse(serviceMonitorTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing: %w", err)
	}

	buf := bytes.Buffer{}

	err = t.Execute(&buf, struct {
		ApplicationName string
		MetricsPath     string
		PortName        string
	}{
		ApplicationName: applicationName,
		MetricsPath:     metricsPath,
		PortName:        config.DefaultApplicationMainPortName,
	})
	if err != nil {
		return nil, fmt.Errorf("executing: %w", err)
	}

	return &buf, nil
}
