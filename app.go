package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/git"
)

type App struct {
	KubeConfigPath    string
	ProjectConfigPath string
	OutputPath        string

	SleepInterval int

	SkipShuffle bool
	SkipFetch   bool
	SkipDeploy  bool
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

	if !app.SkipFetch {
		log.Printf("Wiping output directory '%s'!", app.OutputPath)
		err = os.RemoveAll(app.OutputPath)
		if err != nil {
			return err
		}

		err = os.MkdirAll(app.OutputPath, 0755)
		if err != nil {
			return err
		}
	}

	projectOutputPath := path.Join(app.OutputPath, "config.yml")
	log.Printf("Writing applying configuration to %s", projectOutputPath)
	config.WriteTo(projectOutputPath)
	if err != nil {
		return err
	}

	for i, service := range *config.Services {
		if i != 0 && app.SleepInterval > 0 {
			log.Printf("Sleeping %d seconds...", app.SleepInterval)
			time.Sleep(time.Duration(app.SleepInterval) * time.Second)
		}

		log.Printf("Deploying %s", service.Name)

		if app.SkipFetch {
			log.Printf("Skip fetching manifests via git.")
		} else {
			err := app.FetchService(service)
			if err != nil {
				return err
			}
		}

		if app.SkipDeploy {
			log.Printf("Skip deploying manifests to Kubernetes.")
		} else {
			log.Printf("TODO: kubectl apply")
			if err != nil {
				return err
			}
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
