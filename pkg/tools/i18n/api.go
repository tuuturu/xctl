package i18n

import (
	_ "embed"

	"sigs.k8s.io/yaml"
)

var (
	//go:embed translations/nb_NO/errors.yaml
	rawErrors    []byte            //nolint:gochecknoglobals
	translations map[string]string //nolint:gochecknoglobals
)

func Translate(key string) string {
	if len(translations) == 0 {
		initializeTranslations()
	}

	hit, ok := translations[key]
	if !ok {
		return "!translation not available!"
	}

	return hit
}

func T(key string) string {
	return Translate(key)
}

func initializeTranslations() {
	err := yaml.Unmarshal(rawErrors, &translations)
	if err != nil {
		panic("reading error translations")
	}
}
