package templates

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

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
