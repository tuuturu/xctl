package scaffolding

import (
	_ "embed"
	"io"
	"strings"
)

//go:embed cluster-template.yaml
var clusterTemplate string //nolint:gochecknoglobals

func Cluster() io.Reader {
	return strings.NewReader(clusterTemplate)
}
