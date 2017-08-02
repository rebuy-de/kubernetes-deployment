package api

import (
	"fmt"
	"strings"

	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	log "github.com/sirupsen/logrus"
)

type FetchResult struct {
	Branch    *gh.Branch
	Service   *settings.Service
	Templates map[string]string
}

func (app *App) Fetch(project, branchName string) (*FetchResult, error) {
	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Debugf("fetching templates")

	service := app.Settings.Service(project)
	if service == nil {
		service = app.Settings.GuessService(project)
		log.WithFields(log.Fields{
			"Service": service,
		}).Debug("project not found in settings; guessing it")
	}

	log.WithFields(
		log.Fields(structs.Map(service)),
	).Debug("project found")

	service.Location.Ref = branchName

	branch, err := app.Clients.GitHub.GetBranch(&service.Location)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get branch information")
	}

	log.WithFields(log.Fields{
		"Commit": fmt.Sprintf(
			"commit: %s\n"+
				"Author: %s\n"+
				"Date:   %s\n"+
				"\n%s\n",
			branch.SHA, branch.Author, branch.Date,
			strings.Replace(branch.Message, "\n", "\n    ", -1)),
	}).Debug("fetched latest commit data")

	service.Location.Ref = branch.SHA

	templateStrings, err := app.Clients.GitHub.GetFiles(&service.Location)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &FetchResult{
		Branch:    branch,
		Service:   service,
		Templates: templateStrings,
	}, nil
}
