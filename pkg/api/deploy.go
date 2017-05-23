package api

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

func Deploy(params *Parameters, project, branchName string) {
	settings := params.LoadSettings()

	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Infof("deploying project")

	service := settings.Service(project)
	if service == nil {
		log.WithFields(log.Fields{
			"Project": project,
		}).Fatal("project not found")
		return
	}
	log.WithFields(
		log.Fields(structs.Map(service)),
	).Debug("service found")

	branch, err := params.GitHubClient().GetBranch(&service.Location)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return
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
		log.Fatal(err)
		return
	}

	fmt.Printf("%#v\n", templateStrings)
	fmt.Printf("%#v\n", rendered)
}
