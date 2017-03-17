package cmd

import (
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/imdario/mergo"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

func RenderTemplatesGoal(app *App) error {
	var err error

	err = app.wipeDirectory(renderedSubfolder)
	if err != nil {
		return err
	}

	for _, service := range app.Config.Services {
		manifestInputPath := path.Join(app.Config.Settings.Output, templatesSubfolder, service.Name)
		manifestPath := path.Join(app.Config.Settings.Output, renderedSubfolder, service.Name)
		log.Debugf("Create folder '%s'", manifestPath)

		err := os.MkdirAll(manifestPath, 0755)
		if err != nil {
			return err
		}

		manifests, err := FindFiles(manifestInputPath, "*.yml", "*.yaml")

		log.Debug("Config file template values: %#v", app.Config.Settings.TemplateValues)
		log.Debug("Project template values: %#v", service.TemplateValues)

		mergo.Merge(&service.TemplateValues, app.Config.Settings.TemplateValues)

		log.Debug("Merged template values: %#v", service.TemplateValues)

		for _, manifestInputFile := range manifests {
			err = app.renderTemplate(manifestInputFile, manifestPath, service.TemplateValues)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) renderTemplate(manifestInputFile string, manifestPath string, values map[string]string) error {
	_, manifestFileName := filepath.Split(manifestInputFile)

	manifestOutputFile := path.Join(manifestPath, manifestFileName)
	log.Infof("Templating '%s' to '%s'", manifestInputFile, manifestOutputFile)
	err := templates.ParseManifestFile(manifestInputFile, manifestOutputFile, values)
	if err != nil {
		return err
	}
	return nil
}
