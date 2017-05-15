package cmd

import (
	"io/ioutil"
	"os"
	"path"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/git"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
)

func FetchServicesGoal(app *App) error {
	var err error

	if app.SkipFetch {
		log.Warn("Skip fetching manifests via git.")
		return nil
	}

	err = app.wipeDirectory(templatesSubfolder)
	if err != nil {
		return err
	}

	err = app.wipeDirectory(renderedSubfolder)
	if err != nil {
		return err
	}

	for _, service := range app.Config.Services {
		err := app.FetchService(service)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) FetchService(service *settings.Service) error {
	var err error

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
	service.TemplateValues = append(service.TemplateValues, settings.TemplateValue{
		Name:  "gitCommitID",
		Value: commitID,
	})

	log.Infof("Checked out %s", service.Branch)
	service.TemplateValues = append(service.TemplateValues, settings.TemplateValue{
		Name:  "gitBranchName",
		Value: service.Branch,
	})

	manifests, err := FindFiles(path.Join(tempDir, service.Path), "*.yml", "*.yaml")
	if err != nil {
		return err
	}

	outputPath := path.Join(app.Config.Settings.Output, templatesSubfolder, service.Name)
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
