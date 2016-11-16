package main

import (
	"fmt"
	"os"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/rebuy-de/kubernetes-deployment/kubernetes"
	"github.com/rebuy-de/kubernetes-deployment/settings"
)

type App struct {
	KubectlBuilder    func(kubeconfig *string) (kubernetes.API, error)
	Kubectl           kubernetes.API
	ProjectConfigPath string
	LocalConfigPath   string
	OutputPath        string

	SleepInterval        time.Duration
	IgnoreDeployFailures bool

	RetrySleep time.Duration
	RetryCount int

	SkipShuffle bool
	SkipFetch   bool
	SkipDeploy  bool

	Errors []error
}

const templatesSubfolder = "templates"
const renderedSubfolder = "rendered"

func (app *App) Retry(task Retryer) error {
	return Retry(app.RetryCount, app.RetrySleep, task)
}

func (app *App) Run() error {
	config, err := app.PrepareConfig()
	if err != nil {
		return err
	}

	app.Kubectl, err = app.KubectlBuilder(config.Settings.Kubeconfig)
	if err != nil {
		return err
	}

	err = app.FetchServices(config)
	if err != nil {
		return err
	}

	err = app.RenderTemplates(config)
	if err != nil {
		return err
	}

	err = app.DeployServices(config)
	if err != nil {
		return err
	}

	app.DisplayErrors()

	return nil
}

func (app *App) PrepareConfig() (*settings.ProjectConfig, error) {

	config, err := settings.ReadProjectConfigFrom(app.ProjectConfigPath)
	if err != nil {
		return nil, err
	}

	if app.LocalConfigPath != "" {
		configLoc, err := settings.ReadProjectConfigFrom(app.LocalConfigPath)
		if err != nil {
			return nil, err
		}

		config.MergeConfig(configLoc)
	}

	if err != nil {
		return nil, err
	}

	app.OutputPath = *config.Settings.Output
	app.SleepInterval = *config.Settings.Sleep
	app.RetrySleep = *config.Settings.RetrySleep
	app.RetryCount = *config.Settings.RetryCount
	config.Settings.IgnoreDeployFailures = &app.IgnoreDeployFailures
	config.Settings.SkipShuffle = &app.SkipShuffle
	config.Settings.SkipFetch = &app.SkipFetch
	config.Settings.SkipDeploy = &app.SkipDeploy

	log.Debugf("Read the following project configuration:\n%s", config)

	config.Services.Clean()

	if *config.Settings.SkipShuffle {
		log.Infof("Skip shuffeling service order.")
	} else {
		log.Infof("Shuffling service list")
		config.Services.Shuffle()
	}

	log.Printf("Deploying with this project configuration:\n%s", config)

	if !*config.Settings.SkipFetch {
		log.Warnf("Wiping output directory '%s'!", *config.Settings.Output)
		err := os.RemoveAll(*config.Settings.Output)
		if err != nil {
			return nil, err
		}
	}

	err = os.MkdirAll(*config.Settings.Output, 0755)
	if err != nil {
		return nil, err
	}

	projectOutputPath := path.Join(*config.Settings.Output, "config.yml")
	log.Debugf("Writing applying configuration to %s", projectOutputPath)
	err = config.WriteTo(projectOutputPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (app *App) DisplayErrors() {
	if len(app.Errors) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "\nError(s) occured:\n")
	for i, err := range app.Errors {
		fmt.Fprintf(os.Stderr, "    #%2d: %v\n", i, err)
	}
}
