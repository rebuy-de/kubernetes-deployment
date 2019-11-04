package api

import (
	"fmt"
	"strings"

	argo "github.com/argoproj/argo-cd/pkg/client/clientset/versioned"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func checkForArgoApp(project string) (bool, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return true, errors.Wrap(err, "failed to load kubernetes config for Argo")
	}

	clientset, err := argo.NewForConfig(config)
	if err != nil {
		return true, errors.Wrap(err, "failed to initialize kubernetes client for Argo")
	}

	argoClient := clientset.ArgoprojV1alpha1()

	appsClient := argoClient.Applications("")
	_, err = appsClient.Get(project, v1.GetOptions{})
	if err != nil {
		log.WithFields(log.Fields{
			"Project": project,
			"Error":   err.Error(),
		}).Debug("failed to get argo app - legacy project")
		return false, nil
	}

	return true, errors.New(fmt.Sprintf("Found argo app '%s', abort deployment", project))
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
