package cmd

import (
	"fmt"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubernetes"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
)

type App struct {
	KubectlBuilder func(kubeconfig *string) (kubernetes.API, error)
	Kubectl        kubernetes.API

	Goals  []string
	Config settings.ProjectConfig

	SkipFetch  bool
	SkipDeploy bool
	Target     string
}

const templatesSubfolder = "templates"
const renderedSubfolder = "rendered"

func (app *App) Run() error {
	var err error

	app.Kubectl, err = app.KubectlBuilder(&app.Config.Settings.Kubeconfig)
	if err != nil {
		return err
	}

	goals, err := GetGoals(app.Goals...)
	for _, goal := range goals {
		err = goal(app)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) PrepareConfig() error {
	var err error

	log.Debugf("Read the following configuration:\n%+v", app.Config)

	app.Config.Services.Clean()

	if app.Target != "" {
		services := app.Config.Services
		app.Config.Services = nil

		for _, service := range services {
			if service.Name == app.Target {
				app.Config.Services = settings.Services{
					service,
				}
				break
			}
		}

		if app.Config.Services == nil {
			return fmt.Errorf("Target '%s' not found.", app.Target)
		}
	}

	log.Debugf("Deploying with this project configuration:\n%+v", app.Config)

	err = os.MkdirAll(app.Config.Settings.Output, 0755)
	if err != nil {
		return err
	}

	projectOutputPath := path.Join(app.Config.Settings.Output, "config.yml")
	log.Debugf("Writing applying configuration to %s", projectOutputPath)
	return app.Config.WriteTo(projectOutputPath)
}

func (app *App) wipeDirectory(dir string) error {
	targetDirectory := path.Join(app.Config.Settings.Output, dir)
	log.Infof("Wiping directory %s", targetDirectory)

	err := os.RemoveAll(targetDirectory)
	if err != nil {
		return err
	}

	return os.MkdirAll(targetDirectory, 0755)
}
