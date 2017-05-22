package gh

import (
	"context"
	"path"
	"regexp"
	"time"

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
	GetBranch(location *Location) (*Branch, error)
	GetFile(location *Location) (string, error)
	GetFiles(location *Location) (map[string]string, error)
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

type Branch struct {
	Location
	github.Rate
	Name, SHA string
	Author    string
	Date      time.Time
	Message   string
}

func (gh *API) GetBranch(location *Location) (*Branch, error) {
	ghBranch, resp, err := gh.client.Repositories.GetBranch(
		context.Background(),
		location.Owner, location.Repo, location.Ref,
	)

	if err != nil {
		return nil, err
	}

	branch := &Branch{
		Location: *location,
		Rate:     resp.Rate,
		Name:     *ghBranch.Name,
		Author:   *ghBranch.Commit.Author.Login,
		SHA:      *ghBranch.Commit.SHA,
		Message:  *ghBranch.Commit.Commit.Message,
		Date:     *ghBranch.Commit.Commit.Author.Date,
	}

	log.WithFields(log.Fields{
		"Owner":         location.Owner,
		"Repo":          location.Repo,
		"Ref":           location.Ref,
		"RateLimit":     resp.Rate.Limit,
		"RateRemaining": resp.Rate.Remaining,
		"RateReset":     resp.Rate.Reset,
		"Branch":        *ghBranch.Name,
		"Author":        *ghBranch.Commit.Author.Login,
		"SHA":           *ghBranch.Commit.SHA,
		"Message":       *ghBranch.Commit.Commit.Message,
		"Date":          ghBranch.Commit.Commit.Author.Date,
	}).Debug("fetched branch information")

	return branch, err
}

func (gh *API) GetFile(location *Location) (string, error) {
	log.WithFields(
		log.Fields(structs.Map(location)),
	).Debug("downloading file from GitHub")

	file, _, resp, err := gh.client.Repositories.GetContents(
		context.Background(),
		location.Owner, location.Repo, location.Path,
		&github.RepositoryContentGetOptions{
			Ref: location.Ref,
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

func (gh *API) GetFiles(location *Location) (map[string]string, error) {
	log.WithFields(
		log.Fields(structs.Map(location)),
	).Debug("downloading directory from GitHub")

	_, dir, resp, err := gh.client.Repositories.GetContents(
		context.Background(),
		location.Owner, location.Repo, location.Path,
		&github.RepositoryContentGetOptions{
			Ref: location.Ref,
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
		result[*file.Name], err = gh.GetFile(&Location{
			Owner: location.Owner,
			Repo:  location.Repo,
			Path:  path.Join(location.Path, *file.Name),
			Ref:   location.Ref,
		})
		if err != nil {
			return nil, errors.Wrapf(err,
				"unable to decode file '%v'",
				location)
		}
	}
	return result, nil
}
