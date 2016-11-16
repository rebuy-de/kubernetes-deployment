package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/cloudflare/cfssl/log"
	"github.com/rebuy-de/kubernetes-deployment/settings"
	"github.com/rebuy-de/kubernetes-deployment/templates"
)

func (app *App) RenderTemplates(config *settings.ProjectConfig) error {

	for _, service := range *config.Services {
		manifestInputPath := path.Join(app.OutputPath, templatesSubfolder, service.Name)
		manifestPath := path.Join(app.OutputPath, renderedSubfolder, service.Name)
		log.Debugf("Create folder '%s'", manifestPath)

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
	log.Infof("Templating '%s' to '%s'", manifestInputFile, manifestOutputFile)
	err := templates.ParseManifestFile(manifestInputFile, manifestOutputFile, config.Settings.TemplateValuesMap)
	if err != nil {
		return err
	}
	return nil
}
