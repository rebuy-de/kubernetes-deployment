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
	Templates []gh.File
}

func (app *App) Fetch(project, branchName string) (*FetchResult, error) {
	log.WithFields(log.Fields{
		"Project": project,
		"Branch":  branchName,
	}).Debugf("fetching templates")

	service := app.Settings.Service(project)

	log.WithFields(
		log.Fields(structs.Map(service)),
	).Debug("project found")

	service.Location.Ref = branchName

	isArchived, err := app.GitHub.IsArchived(&service.Location)
	if err != nil {
		return nil, errors.Wrap(err, "unable to check for archived repo")
	}
	if isArchived {
		return nil, errors.Wrap(err, "repo is archived, please use active repo")
	}

	branch, err := app.GitHub.GetBranch(&service.Location)
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

	files, err := app.GitHub.GetFiles(&service.Location)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	app.StartInterceptors(service)

	err = app.Interceptors.PostFetch(branch)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &FetchResult{
		Branch:    branch,
		Service:   service,
		Templates: files,
	}, nil
}
