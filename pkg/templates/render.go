package templates

import (
	"bytes"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func RenderAll(templates []gh.File, variables Variables) ([]gh.File, error) {
	var result []gh.File
	for _, file := range templates {
		rendered, err := Render(file.Content, variables)
		if err != nil {
			return nil, errors.Wrapf(err, "Unable to render '%s'", file.Name)
		}
		result = append(result, gh.File{Path: file.Path, Content: rendered})
	}
	return result, nil
}

func Render(templateString string, variables Variables) (string, error) {
	funcMap := template.FuncMap{
		"ToUpper":    strings.ToUpper,
		"ToLower":    strings.ToLower,
		"Identifier": IdentifierFunc,
		"MakeSlice":  MakeSliceFunc,
	}

	t, err := template.
		New("").
		Funcs(funcMap).
		Parse(templateString)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to parse template")
	}

	var buf bytes.Buffer

	err = t.Execute(&buf, variables)
	if err != nil {
		return "", errors.Wrapf(err, "Unable to render template")
	}

	log.WithFields(log.Fields{
		"Result": buf.String(),
	}).Debug("Rendered file")

	return buf.String(), nil
}
