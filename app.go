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

	Commands []Command
	Config   *settings.ProjectConfig

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
	err := app.PrepareConfig()
	if err != nil {
		return err
	}

	app.Kubectl, err = app.KubectlBuilder(app.Config.Settings.Kubeconfig)
	if err != nil {
		return err
	}

	for _, command := range app.Commands {
		err = command(app)
		if err != nil {
			return err
		}
	}

	app.DisplayErrors()

	return nil
}

func (app *App) PrepareConfig() error {
	config, err := settings.ReadProjectConfigFrom(app.ProjectConfigPath)
	if err != nil {
		return err
	}

	if app.LocalConfigPath != "" {
		configLoc, err := settings.ReadProjectConfigFrom(app.LocalConfigPath)
		if err != nil {
			return err
		}

		config.MergeConfig(configLoc)
	}

	if err != nil {
		return err
	}

	app.OutputPath = *config.Settings.Output
	app.SleepInterval = *config.Settings.Sleep
	app.RetrySleep = *config.Settings.RetrySleep
	app.RetryCount = *config.Settings.RetryCount
	config.Settings.IgnoreDeployFailures = &app.IgnoreDeployFailures

	log.Debugf("Read the following project configuration:\n%s", config)

	config.Services.Clean()

	fmt.Printf("%#v\n", *config.Settings)

	if config.Settings.SkipShuffle != nil && *config.Settings.SkipShuffle {
		log.Infof("Skip shuffeling service order.")
	} else {
		log.Infof("Shuffling service list")
		config.Services.Shuffle()
	}

	log.Printf("Deploying with this project configuration:\n%s", config)

	log.Warnf("Wiping output directory '%s'!", *config.Settings.Output)
	err = os.RemoveAll(*config.Settings.Output)
	if err != nil {
		return err
	}

	err = os.MkdirAll(*config.Settings.Output, 0755)
	if err != nil {
		return err
	}

	projectOutputPath := path.Join(*config.Settings.Output, "config.yml")
	log.Debugf("Writing applying configuration to %s", projectOutputPath)
	err = config.WriteTo(projectOutputPath)
	if err != nil {
		return err
	}

	app.Config = config
	return nil
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
