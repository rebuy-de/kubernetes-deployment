package gh

import (
	"context"
	"path"
	"regexp"

	"golang.org/x/oauth2"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/google/go-github/github"
	"github.com/pkg/errors"
)

var (
	ContentLocationRE = regexp.MustCompile(`^github.com/([^/]+)/([^/]+)/(.*)$`)
)

type Client interface {
	GetFile(location Location) (string, error)
	GetFiles(location Location) (map[string]string, error)
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

func (gh *API) GetFile(location Location) (string, error) {
	log.WithFields(
		log.Fields(structs.Map(location)),
	).Debug("downloading file from GitHub")

	file, _, resp, err := gh.client.Repositories.GetContents(
		context.Background(),
		location.Owner, location.Repo, location.Path,
		&github.RepositoryContentGetOptions{
			Ref: location.Branch,
		},
	)

	if err != nil {
		return "", errors.Wrapf(err,
			"unable to fetch file '%v' from GitHub", location)
	}

	if file == nil {
		return "", errors.Errorf(
			"unable to fetch file '%v' from GitHub; probably it's a directoy",
			location)
	}

	log.WithFields(log.Fields{
		"Size":          *file.Size,
		"URL":           *file.HTMLURL,
		"RateLimit":     resp.Rate.Limit,
		"RateRemaining": resp.Rate.Remaining,
		"RateReset":     resp.Rate.Reset,
	}).Debug("found file")

	return file.GetContent()
}

func (gh *API) GetFiles(location Location) (map[string]string, error) {
	log.WithFields(
		log.Fields(structs.Map(location)),
	).Debug("downloading directory from GitHub")

	_, dir, resp, err := gh.client.Repositories.GetContents(
		context.Background(),
		location.Owner, location.Repo, location.Path,
		&github.RepositoryContentGetOptions{
			Ref: location.Branch,
		},
	)

	if err != nil {
		return nil, errors.Wrapf(err,
			"unable to fetch directory '%v' from GitHub", location)
	}

	if dir == nil {
		return nil, errors.Errorf(
			"unable to fetch directory '%v' from GitHub; probably it's a file",
			location)
	}

	log.WithFields(log.Fields{
		"Files":         len(dir),
		"RateLimit":     resp.Rate.Limit,
		"RateRemaining": resp.Rate.Remaining,
		"RateReset":     resp.Rate.Reset,
	}).Debug("found files in directory")

	result := make(map[string]string)
	for _, file := range dir {
		result[*file.Name], err = gh.GetFile(Location{
			Owner:  location.Owner,
			Repo:   location.Repo,
			Path:   path.Join(location.Path, *file.Name),
			Branch: location.Branch,
		})
		if err != nil {
			return nil, errors.Wrapf(err,
				"unable to decode file '%v'",
				location)
		}
	}
	return result, nil
}
