package api

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api"
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

	decode := api.Codecs.UniversalDeserializer().Decode

	objects := []runtime.Object{}

	for name, data := range rendered {
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			log.WithFields(log.Fields{
				"Name": name,
			}).Debug("Ignoring file with wrong extension.")
			continue
		}

		if strings.TrimSpace(data) == "" {
			log.WithFields(log.Fields{
				"Name": name,
			}).Debug("Ignoring empty file.")
			continue
		}

		obj, _, err := decode([]byte(data), nil, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode file '%s'", name)
		}

		obj, err = app.Interceptors.ManifestRendered(obj)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		objects = append(objects, obj)
	}

	return objects, nil
}
