package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rebuy-de/kubernetes-deployment/git"
)

type App struct {
	KubeConfigPath string

	ProjectConfigPath string
	ProjectOutputPath string

	CheckoutDirectory string

	SleepInterval int
	SkipShuffle   bool
}

func (app *App) Run() error {
	config, err := ReadProjectConfigFrom(app.ProjectConfigPath)
	if err != nil {
		return err
	}

	log.Printf("Read the following project configuration:\n%s", config)

	config.Services.Clean()

	if !app.SkipShuffle {
		log.Printf("Shuffling service list")
		config.Services.Shuffle()
	}

	log.Printf("Deploying with this project configuration:\n%s", config)

	log.Printf("Writing applying configuration to %s", app.ProjectOutputPath)
	config.WriteTo(app.ProjectOutputPath)
	if err != nil {
		return err
	}

	for _, service := range *config.Services {
		app.DeployService(service)
	}

	return nil
}

func (app *App) DeployService(service *Service) error {
	log.Printf("Deploying %+v", *service)

	tempDir, err := ioutil.TempDir("", "kubernetes-deployment-checkout-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	err = git.SparseCheckout(tempDir, service.Repository, service.Branch, service.Path)
	if err != nil {
		return err
	}

	manifests, err := findManifests(path.Join(tempDir, service.Path))
	if err != nil {
		return err
	}
	log.Printf("Found these manifests: %v", manifests)

	return nil
}

func findManifests(dir string) ([]string, error) {
	dir = path.Clean(dir) + "/"
	result := []string{}

	for _, ext := range []string{"*.yml", "*.yaml"} {
		matches, err := filepath.Glob(path.Join(dir, ext))
		if err != nil {
			return nil, err
		}

		for _, m := range matches {
			m = strings.TrimPrefix(m, dir)
			result = append(result, m)
		}
	}

	return result, nil
}
