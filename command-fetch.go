package main

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/rebuy-de/kubernetes-deployment/git"
	"github.com/rebuy-de/kubernetes-deployment/settings"
)

func (app *App) FetchServices(config *settings.ProjectConfig) error {
	if app.SkipFetch {
		log.Warn("Skip fetching manifests via git.")
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
		log.Infof("Copying manifest to '%s'", target)

		err := CopyFile(manifest, target)
		if err != nil {
			return err
		}
	}
	return nil
}
