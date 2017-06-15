package templates

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

func RenderAll(templates map[string]string, values Values) (map[string]string, error) {
	result := map[string]string{}
	for name, templateString := range templates {
		rendered, err := Render(templateString, values)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to render '%s'", name)
		}
		result[name] = rendered
	}
	return result, nil
}

func Render(templateString string, values Values) (string, error) {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	t, err := template.
		New("").
		Funcs(funcMap).
		Parse(templateString)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to parse template")
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, values)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to render template")
	}

	return buf.String(), nil
}
