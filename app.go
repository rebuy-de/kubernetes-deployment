package main

import (
	"fmt"
	"github.com/rebuy-de/kubernetes-deployment/git"
	"github.com/rebuy-de/kubernetes-deployment/kubernetes"
	"github.com/rebuy-de/kubernetes-deployment/templates"
	"github.com/rebuy-de/kubernetes-deployment/settings"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

type App struct {
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

	err = app.FetchServices(config)
	if err != nil {
		return err
	}

	err = app.RenderTemplates(config)
	if err != nil {
		return err
	}

	if app.Kubectl == nil {
		app.Kubectl, err = kubernetes.New(*config.Settings.Kubeconfig)
		if err != nil {
			return err
		}
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
	app.IgnoreDeployFailures = *config.Settings.IgnoreDeployFailures
	app.RetrySleep = *config.Settings.RetrySleep
	app.RetryCount = *config.Settings.RetryCount
	app.SkipShuffle = *config.Settings.SkipShuffle
	app.SkipFetch = *config.Settings.SkipFetch
	app.SkipDeploy = *config.Settings.SkipDeploy



	log.Printf("Read the following project configuration:\n%s", config)

	config.Services.Clean()

	if *config.Settings.SkipShuffle {
		log.Printf("Skip shuffeling service order.")
	} else {
		log.Printf("Shuffling service list")
		config.Services.Shuffle()
	}

	log.Printf("Deploying with this project configuration:\n%s", config)

	if !*config.Settings.SkipFetch {
		log.Printf("Wiping output directory '%s'!", *config.Settings.Output)
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
	log.Printf("Writing applying configuration to %s", projectOutputPath)
	err = config.WriteTo(projectOutputPath)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (app *App) FetchServices(config *settings.ProjectConfig) error {
	if app.SkipFetch {
		log.Printf("Skip fetching manifests via git.")
		return nil
	}

	for _, service := range *config.Services {
		err := app.Retry(func() error {
			return app.FetchService(service, config)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) FetchService(service *settings.Service, config *settings.ProjectConfig) error {
	tempDir, err := ioutil.TempDir("", "kubernetes-deployment-checkout-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = git.SparseCheckout(tempDir, service.Repository, service.Branch, service.Path)
	if err != nil {
		return err
	}

	manifests, err := FindFiles(path.Join(tempDir, service.Path), "*.yml", "*.yaml")
	if err != nil {
		return err
	}

	outputPath := path.Join(app.OutputPath, templatesSubfolder, service.Name)
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		return err
	}

	for _, manifest := range manifests {
		name := path.Base(manifest)
		target := path.Join(outputPath, name)
		log.Printf("Copying manifest to '%s'", target)

		err := CopyFile(manifest, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *App) RenderTemplates(config *settings.ProjectConfig) error {

	for _, service := range *config.Services {
		manifestInputPath := path.Join(app.OutputPath, templatesSubfolder, service.Name)
		manifestPath := path.Join(app.OutputPath, renderedSubfolder, service.Name)
		log.Printf("Create folder '%s'", manifestPath)

		err := os.MkdirAll(manifestPath, 0755)
		if err != nil {
			return err
		}

		manifests, err := FindFiles(manifestInputPath, "*.yml", "*.yaml")

		for _, manifestInputFile := range manifests {
			err = app.renderTemplate(manifestInputFile, manifestPath, config)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) renderTemplate(manifestInputFile string, manifestPath string, config *settings.ProjectConfig) error {
	_, manifestFileName := filepath.Split(manifestInputFile)

	manifestOutputFile := path.Join(manifestPath, manifestFileName)
	log.Printf("Templating '%s' to '%s'", manifestInputFile, manifestOutputFile)
	err := templates.ParseManifestFile(manifestInputFile, manifestOutputFile, config.Settings.TemplateValuesMap)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) DeployServices(config *settings.ProjectConfig) error {
	for i, service := range *config.Services {
		if i != 0 && app.SleepInterval > 0 {
			log.Printf("Sleeping %v ...", app.SleepInterval)
			time.Sleep(app.SleepInterval)
		}

		if app.SkipDeploy {
			log.Printf("Skip deploying manifests to Kubernetes.")
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
	manifestPath := path.Join(app.OutputPath, renderedSubfolder, service.Name)
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

		log.Printf("Applying manifest '%s'", manifestInputFile)
		err := app.Retry(func() error {
			_, err := app.Kubectl.Apply(manifestInputFile)
			return err
		})
		if err != nil && app.IgnoreDeployFailures {
			log.Printf("Ignoring failed deployment of %s", service.Name)
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

func (app *App) DisplayErrors() {

	fmt.Fprintf(os.Stderr, "\nError(s) occured:\n")
	for i, err := range app.Errors {
		fmt.Fprintf(os.Stderr, "    #%2d: %v\n", i, err)
	}
}
