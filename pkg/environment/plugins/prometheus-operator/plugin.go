package prometheus_operator

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/deifyed/xctl/pkg/apis/xctl/v1alpha1"
	"sigs.k8s.io/yaml"
)

func plugin() v1alpha1.Plugin {
	t := template.Must(template.New("plugin").Parse(pluginTemplate))

	buf := bytes.Buffer{}

	err := t.Execute(&buf, struct{}{})
	if err != nil {
		panic(err)
	}

	var plugin v1alpha1.Plugin

	err = yaml.Unmarshal(buf.Bytes(), &plugin)
	if err != nil {
		panic(err)
	}

	return plugin
}

//go:embed plugin.yaml
var pluginTemplate string
