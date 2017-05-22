package api

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
)

func Deploy(params *Parameters, project, branch string) {
	settings := params.LoadSettings()

	log.Infof("deploying %s/%s", project, branch)

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

	templates, err := params.GitHubClient().GetFiles(service.Location)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("%#v\n", templates)
}
