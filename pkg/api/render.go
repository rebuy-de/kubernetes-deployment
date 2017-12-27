package api

import (
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func (app *App) Render(fetched *FetchResult) ([]runtime.Object, error) {
	defaultVariables := templates.Variables{
		"gitBranchName": fetched.Branch.Name,
		"gitCommitID":   fetched.Branch.SHA,
	}

	fetched.Service.Variables.Defaults(defaultVariables)

	log.WithFields(log.Fields{
		"Values": fetched.Service.Variables,
	}).Debug("collected template values")

	rendered, err := templates.RenderAll(fetched.Templates, fetched.Service.Variables)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return app.decode(rendered)
}

func (app *App) decode(rendered map[string]string) ([]runtime.Object, error) {
	var objects []runtime.Object
	decode := scheme.Codecs.UniversalDeserializer().Decode

	for name, data := range rendered {
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			log.WithFields(log.Fields{
				"Name": name,
			}).Debug("Ignoring file with wrong extension.")
			continue
		}

		splitted := regexp.MustCompile("[\n\r]---").Split(data, -1)

		for _, part := range splitted {
			if strings.TrimSpace(part) == "" {
				log.WithFields(log.Fields{
					"Name": name,
				}).Debug("Ignoring empty file.")
				continue
			}

			obj, _, err := decode([]byte(part), nil, nil)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to decode file '%s'", name)
			}

			obj, err = app.Interceptors.PostManifestRender(obj)
			if err != nil {
				return nil, errors.WithStack(err)
			}

			objects = append(objects, obj)
		}
	}

	return objects, nil
}
