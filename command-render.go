package main

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/settings"
	"github.com/rebuy-de/kubernetes-deployment/templates"
)

func RenderTemplatesCommand(app *App) error {
	var err error

	err = app.wipeDirectory(renderedSubfolder)
	if err != nil {
		return err
	}

	for _, service := range *app.Config.Services {
		manifestInputPath := path.Join(app.OutputPath, templatesSubfolder, service.Name)
		manifestPath := path.Join(app.OutputPath, renderedSubfolder, service.Name)
		log.Debugf("Create folder '%s'", manifestPath)

		err := os.MkdirAll(manifestPath, 0755)
		if err != nil {
			return err
		}

		manifests, err := FindFiles(manifestInputPath, "*.yml", "*.yaml")

		for _, manifestInputFile := range manifests {
			err = app.renderTemplate(manifestInputFile, manifestPath, app.Config.Settings.TemplateValues.Merge(service.TemplateValues))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) renderTemplate(manifestInputFile string, manifestPath string, values settings.TemplateValues) error {
	_, manifestFileName := filepath.Split(manifestInputFile)

	manifestOutputFile := path.Join(manifestPath, manifestFileName)
	log.Infof("Templating '%s' to '%s'", manifestInputFile, manifestOutputFile)
	err := templates.ParseManifestFile(manifestInputFile, manifestOutputFile, values)
	if err != nil {
		return err
	}
	return nil
}
