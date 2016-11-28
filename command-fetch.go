package main

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/git"
	"github.com/rebuy-de/kubernetes-deployment/settings"
)

func FetchServicesCommand(app *App) error {
	if app.SkipFetch {
		log.Warn("Skip fetching manifests via git.")
		return nil
	}

	for _, service := range *app.Config.Services {
		err := app.Retry(func() error {
			return app.FetchService(service, app.Config)
		})
		if err != nil {
			return err
		}
	}

	projectOutputPath := path.Join(*app.Config.Settings.Output, "config.yml")
	log.Debugf("Writing applying configuration to %s", projectOutputPath)
	return app.Config.WriteTo(projectOutputPath)
}

func (app *App) FetchService(service *settings.Service, config *settings.ProjectConfig) error {
	tempDir, err := ioutil.TempDir("", "kubernetes-deployment-checkout-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	repo, err := git.SparseCheckout(tempDir, service.Repository, service.Branch, service.Path)
	if err != nil {
		return err
	}

	commitID, err := repo.CommitID()
	if err != nil {
		return err
	}

	log.Infof("Checked out %s", commitID)
	service.TemplateValues = service.TemplateValues.Merge(settings.TemplateValues{
		"gitCommitID": commitID,
	})

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
