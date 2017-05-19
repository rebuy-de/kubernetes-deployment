package gh

import (
	"context"
	"regexp"

	"golang.org/x/oauth2"

	log "github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

var (
	ContentLocationRE = regexp.MustCompile(`^github.com/([^/]+)/([^/]+)/(.*)$`)
)

type Client interface {
	GetContents(location string) (string, error)
}

type API struct {
	client *github.Client
}

func New(token string) Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &API{
		client: client,
	}
}

func (gh *API) GetContents(location string) (string, error) {
	matches := ContentLocationRE.FindStringSubmatch(location)
	if matches == nil {
		return "", errors.Errorf(
			"GitHub location must have the form `github.com/:owner:/:repo:/:path:`")
	}

	var (
		owner = matches[1]
		repo  = matches[2]
		path  = matches[3]
	)

	log.WithFields(log.Fields{
		"Owner": owner,
		"Repo":  repo,
		"Path":  path,
	}).Debug("downloading file from GitHub")

	file, _, _, err := gh.client.Repositories.GetContents(
		context.Background(),
		owner, repo, path,
		&github.RepositoryContentGetOptions{
			Ref: "master",
		},
	)

	if err != nil {
		return "", errors.Wrapf(err,
			"unable to fetch file '%s' from GitHub", location)
	}

	if file == nil {
		return "", errors.Errorf(
			"unable to fetch file '%s' from GitHub; probably it's a directoy",
			location)
	}

	return file.GetContent()
}
