package api

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
	log "github.com/sirupsen/logrus"
)

func (app *App) Generate(project, branchName string) ([]runtime.Object, error) {
	app.Settings.Clean(app.Parameters.Context)

	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Infof("deploying project")

	service := app.Settings.Service(project)
	if service == nil {
		return nil, errors.Errorf("project '%s' not found", project)
	}

	log.WithFields(
		log.Fields(structs.Map(service)),
	).Debug("project found")

	service.Location.Ref = branchName

	branch, err := app.Clients.GitHub.GetBranch(&service.Location)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to get branch information")
	}

	log.Infof("latest commit:\n\n"+
		"commit: %s\n"+
		"Author: %s\n"+
		"Date:   %s\n"+
		"\n%s\n",
		branch.SHA, branch.Author, branch.Date,
		strings.Replace(branch.Message, "\n", "\n    ", -1))

	service.Location.Ref = branch.SHA

	templateStrings, err := app.Clients.GitHub.GetFiles(&service.Location)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defaultVariables := templates.Variables{
		"gitBranchName": branch.Name,
		"gitCommitID":   branch.SHA,
	}

	service.Variables.Defaults(defaultVariables)

	log.WithFields(log.Fields{
		"Values": service.Variables,
	}).Debug("collected template values")

	rendered, err := templates.RenderAll(templateStrings, service.Variables)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	decode := api.Codecs.UniversalDeserializer().Decode

	objects := []runtime.Object{}

	for name, data := range rendered {
		if !strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml") {
			log.WithFields(log.Fields{
				"Name": name,
			}).Debug("Ignoring file with wrong extension.")
			continue
		}

		obj, _, err := decode([]byte(data), nil, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to decode file '%s'", name)
		}

		objects = append(objects, obj)
	}

	if len(rendered) <= 0 {
		return nil, errors.Errorf("directory doesn't contain any template files")
	}

	return objects, nil
}
