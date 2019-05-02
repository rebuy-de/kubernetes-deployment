package api

import (
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	log "github.com/sirupsen/logrus"
)

func (app *App) Generate(project, branchName string) ([]runtime.Object, error) {
	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Debugf("generating manifests")

	app.Statsd.Increment("generate",
		statsdw.Tag{Name: "project", Value: project},
		statsdw.Tag{Name: "branch", Value: branchName})

	fetched, err := app.Fetch(project, branchName)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch project")
	}

	objects, err := app.Render(fetched)
	if err != nil {
		return nil, errors.Wrap(err, "failed to render project")
	}

	if len(objects) <= 0 {
		return nil, errors.Errorf("didn't find any template files")
	}

	return objects, nil
}
