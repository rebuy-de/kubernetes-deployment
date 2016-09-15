package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
	"github.com/rebuy-de/kubernetes-deployment/git"
	"github.com/rebuy-de/kubernetes-deployment/kubernetes"
)

type App struct {
	Kubectl              kubernetes.API
	ProjectConfigPath    string
	LocalConfigPath      string
	OutputPath           string

	SleepInterval        time.Duration
	IgnoreDeployFailures bool

	RetrySleep           time.Duration
	RetryCount           int

	SkipShuffle          bool
	SkipFetch            bool
	SkipDeploy           bool

	Errors               []error
}

func (app *App) Retry(task Retryer) error {
	return Retry(app.RetryCount, app.RetrySleep, task)
}

func (app *App) Run() error {
	config, err := app.PrepareConfig()
	if err != nil {
		return err
	}

	app.OutputPath = *config.Settings.Output
	app.SleepInterval = *config.Settings.Sleep
	app.IgnoreDeployFailures = *config.Settings.IgnoreDeployFailures
	app.RetrySleep = *config.Settings.RetrySleep
	app.RetryCount = *config.Settings.RetryCount
	app.SkipShuffle = *config.Settings.SkipShuffle
	app.SkipFetch = *config.Settings.SkipFetch
	app.SkipDeploy = *config.Settings.SkipDeploy

	app.Kubectl, err = kubernetes.New(*config.Settings.Kubeconfig)
	if err != nil {
		return err
	}

	err = app.FetchServices(config)
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

func (app *App) PrepareConfig() (*ProjectConfig, error) {

	config, err := ReadProjectConfigFrom(app.ProjectConfigPath)
	if err != nil {
		return nil, err
	}

	configLoc, err := ReadProjectConfigFrom(app.LocalConfigPath)
	if err != nil {
		return nil, err
	}

	mergeConfig(config, configLoc)


	log.Printf("Read the following project configuration:\n%s", config)

	config.Services.Clean()

	if app.SkipShuffle {
		log.Printf("Skip shuffeling service order.")
	} else {
		log.Printf("Shuffling service list")
		config.Services.Shuffle()
	}

	log.Printf("Deploying with this project configuration:\n%s", config)

	if !app.SkipFetch {
		log.Printf("Wiping output directory '%s'!", app.OutputPath)
		err := os.RemoveAll(app.OutputPath)
		if err != nil {
			return nil, err
		}
	}

	err = os.MkdirAll(app.OutputPath, 0755)
	if err != nil {
		return nil, err
	}

	projectOutputPath := path.Join(app.OutputPath, "config.yml")
	log.Printf("Writing applying configuration to %s", projectOutputPath)
	err = config.WriteTo(projectOutputPath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func mergeConfig(defaultConfig *ProjectConfig, localConfig *ProjectConfig) {

	if localConfig.Settings.Kubeconfig != nil {
		defaultConfig.Settings.Kubeconfig = localConfig.Settings.Kubeconfig
	}

	if localConfig.Settings.Output != nil {
		defaultConfig.Settings.Output = localConfig.Settings.Output
	}

	if localConfig.Settings.Sleep != nil {
		defaultConfig.Settings.Sleep = localConfig.Settings.Sleep
	}

	if localConfig.Settings.SkipShuffle != nil {
		defaultConfig.Settings.SkipShuffle = localConfig.Settings.SkipShuffle
	}

	if localConfig.Settings.SkipFetch != nil {
		defaultConfig.Settings.SkipFetch = localConfig.Settings.SkipFetch
	}

	if localConfig.Settings.SkipDeploy != nil {
		defaultConfig.Settings.SkipDeploy = localConfig.Settings.SkipDeploy
	}

	if localConfig.Settings.RetrySleep != nil {
		defaultConfig.Settings.RetrySleep = localConfig.Settings.RetrySleep
	}

	if localConfig.Settings.RetryCount != nil {
		defaultConfig.Settings.RetryCount = localConfig.Settings.RetryCount
	}

	if localConfig.Settings.IgnoreDeployFailures != nil {
		defaultConfig.Settings.IgnoreDeployFailures = localConfig.Settings.IgnoreDeployFailures
	}

	tempMap := make(map[string]string)

	for _, templateValue := range *defaultConfig.Settings.TemplateValues {
		tempMap[templateValue.Key] = templateValue.Value
	}

	for _, templateValue := range *localConfig.Settings.TemplateValues {
		tempMap[templateValue.Key] = templateValue.Value
	}

	defaultConfig.Settings.TemplateValuesMap = tempMap

	fmt.Println(defaultConfig)

	os.Exit(1)
}


func (app *App) FetchServices(config *ProjectConfig) error {
	if app.SkipFetch {
		log.Printf("Skip fetching manifests via git.")
		return nil
	}

	for _, service := range *config.Services {
		err := app.Retry(func() error {
			return app.FetchService(service)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) FetchService(service *Service) error {
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

	outputPath := path.Join(app.OutputPath, service.Name)
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

func (app *App) DeployServices(config *ProjectConfig) error {
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

func (app *App) DeployService(service *Service) error {
	manifestPath := path.Join(app.OutputPath, service.Name)
	manifests, err := FindFiles(manifestPath, "*.yml", "*.yaml")
	if err != nil {
		return err
	}

	if len(manifests) == 0 {
		return fmt.Errorf("Did not find any manifest for '%s' in '%s'",
			service.Name, manifestPath)
	}

	for _, manifest := range manifests {
		log.Printf("Applying manifest '%s'", manifest)
		err := app.Retry(func() error {
			_, err := app.Kubectl.Apply(manifest)
			return err
		})
		if err != nil && app.IgnoreDeployFailures {
			log.Printf("Ignoring failed deployment of %s", service.Name)
			app.Errors = append(app.Errors,
				fmt.Errorf("Deployment of '%s' in service '%s' failed: %v",
					manifest, service.Name, err),
			)
		}
		if err != nil && !app.IgnoreDeployFailures {
			return err
		}
	}

	return nil
}

func (app *App) DisplayErrors() {
	if len(app.Errors) < 1 {
		return
	}

	fmt.Fprintf(os.Stderr, "\nError(s) occured:\n")
	for i, err := range app.Errors {
		fmt.Fprintf(os.Stderr, "    #%2d: %v\n", i, err)
	}
}
