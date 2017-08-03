package api

import (
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	log "github.com/sirupsen/logrus"
)

func (app *App) Apply(project, branchName string) error {
	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Debug("applying manifests")

	app.Clients.Statsd.Increment("apply",
		statsdw.Tag{Name: "project", Value: project},
		statsdw.Tag{Name: "branch", Value: branchName})

	objects, err := app.Generate(project, branchName)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, obj := range objects {
		upstreamObj, err := app.Clients.Kubectl.Apply(obj)
		if err != nil {
			return errors.Wrap(err, "unable to apply manifest")
		}
		err = app.Interceptors.PreManifestApply(upstreamObj)
		if err != nil {
			return errors.WithStack(err)
		}

		log.WithFields(log.Fields{
			"Manifest": upstreamObj,
		}).Debug("applied manifest")
	}

	err = app.Interceptors.PostApply(objects)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
