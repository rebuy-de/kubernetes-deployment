package api

import (
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

func Generate(params *Parameters, project, branchName string) ([]runtime.Object, error) {
	settings := params.LoadSettings()

	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Infof("deploying project")

	service := settings.Service(project)
	if service == nil {
		return nil, errors.Errorf("project '%s' not found", project)
	}

	log.WithFields(
		log.Fields(structs.Map(service)),
	).Debug("project found")

	branch, err := params.GitHubClient().GetBranch(&service.Location)
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

	templateStrings, err := params.GitHubClient().GetFiles(&service.Location)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defaultValues := templates.Values{
		"gitBranchName": branch.Name,
		"gitCommitID":   branch.SHA,
	}

	service.TemplateValues.Defaults(defaultValues)

	log.WithFields(log.Fields{
		"Values": service.TemplateValues,
	}).Debug("collected template values")

	rendered, err := templates.RenderAll(templateStrings, service.TemplateValues)
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

	return objects, nil
}
