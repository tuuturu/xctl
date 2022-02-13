package i18n

import (
	_ "embed"

	"sigs.k8s.io/yaml"
)

var (
	//go:embed translations/nb_NO/errors.yaml
	rawErrors []byte //nolint:gochecknoglobals
	//go:embed translations/nb_NO/cmd.yaml
	rawCmd       []byte            //nolint:gochecknoglobals
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
	rawTranslations := [][]byte{
		rawErrors,
		rawCmd,
	}

	marshalledTranslations := make([]map[string]string, 0)

	for _, rawTranslation := range rawTranslations {
		currentTranslations := make(map[string]string)

		err := yaml.Unmarshal(rawTranslation, &currentTranslations)
		if err != nil {
			panic("reading error translations")
		}

		marshalledTranslations = append(marshalledTranslations, currentTranslations)
	}

	translations = mergeTranslations(marshalledTranslations...)
}

func mergeTranslations(args ...map[string]string) map[string]string {
	result := make(map[string]string)

	for _, translation := range args {
		for key, value := range translation {
			result[key] = value
		}
	}

	return result
}
