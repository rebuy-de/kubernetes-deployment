package api

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
)

func Deploy(params *Parameters, project, branchName string) {
	settings := params.LoadSettings()

	log.Infof("deploying %s/%s", project, branchName)

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

	service.Location.Ref = branch.SHA

	templates, err := params.GitHubClient().GetFiles(&service.Location)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%#v\n", branch)
	fmt.Printf("%#v\n", templates)
}
