package cmd

import (
	"fmt"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
)

func DeployServicesGoal(app *App) error {
	for _, service := range app.Config.Services {
		if app.SkipDeploy {
			log.Warn("Skip deploying manifests to Kubernetes.")
		} else {
			err := app.DeployService(service)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) DeployService(service *settings.Service) error {
	manifestPath := path.Join(app.Config.Settings.Output, renderedSubfolder, service.Name)
	manifests, err := FindFiles(manifestPath, "*.yml", "*.yaml")
	if err != nil {
		return err
	}

	if len(manifests) == 0 {
		return fmt.Errorf("Did not find any manifest for '%s' in '%s'",
			service.Name, manifestPath)
	}

	for _, manifestInputFile := range manifests {
		if err != nil {
			return err
		}

		log.Infof("Applying manifest '%s'", manifestInputFile)
		_, err := app.Kubectl.Apply(manifestInputFile)
		if err != nil && app.IgnoreDeployFailures {
			log.Errorf("Ignoring failed deployment of %s", service.Name)
			app.Errors = append(app.Errors,
				fmt.Errorf("Deployment of '%s' in service '%s' failed: %v",
					manifestInputFile, service.Name, err),
			)
		}
		if err != nil && !app.IgnoreDeployFailures {
			return err
		}
	}
	return nil
}
