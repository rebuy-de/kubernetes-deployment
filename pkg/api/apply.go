package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	log "github.com/sirupsen/logrus"
)

func checkForArgoApp(project string) (bool, error) {
	req, _ := http.NewRequest("GET", "https://argocd.production.rebuy.cloud/api/v1/applications/"+project, nil)
	argoToken := os.Getenv("ARGOCD_API_TOKEN")
	if argoToken == "" {
		log.WithFields(log.Fields{
			"Project": project,
		}).Debug("No ArgoCD API token found - continuing deployment")
		return false, nil
	}
	req.AddCookie(&http.Cookie{Name: "argocd.token", Value: argoToken})
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return true, errors.Wrap(err, "Unable to query ArgoCD API for application")
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.WithFields(log.Fields{
			"Project": project,
		}).Debug("failed to get argo app - legacy project")
		return false, nil
	}

	return true, errors.New(fmt.Sprintf("Project is managed by ArgoCD. Please deploy with `/kubot-deploy %s`.", project))
}

func (app *App) Apply(project, branchName string) error {
	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Debug("applying manifests")

	app.Statsd.Increment("apply",
		statsdw.Tag{Name: "project", Value: project},
		statsdw.Tag{Name: "branch", Value: branchName})

	isArgoProject, err := checkForArgoApp(strings.Split(project, "/")[0])
	if err != nil && isArgoProject == true {
		return errors.WithStack(err)
	}

	objects, err := app.Generate(project, branchName)
	if err != nil {
		return errors.WithStack(err)
	}

	err = app.Interceptors.PreApply(objects)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, obj := range objects {
		obj, err = app.Interceptors.PreManifestApply(obj)
		if err != nil {
			return errors.WithStack(err)
		}

		upstreamObj, err := app.Kubectl.Apply(obj)
		if err != nil {
			return errors.Wrap(err, "unable to apply manifest")
		}

		err = app.Interceptors.PostManifestApply(upstreamObj)
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
