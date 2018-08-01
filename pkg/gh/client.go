package gh

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"regexp"
	"strings"
	"time"

	"golang.org/x/oauth2"

	"github.com/fatih/structs"
	"github.com/google/go-github/github"
	"github.com/gregjones/httpcache"
	"github.com/gregjones/httpcache/diskcache"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	log "github.com/sirupsen/logrus"
)

var (
	ContentLocationRE = regexp.MustCompile(`^github.com/([^/]+)/([^/]+)/(.*)$`)
)

type Interface interface {
	GetBranch(location *Location) (*Branch, error)
	GetFile(location *Location) (File, error)
	GetFiles(location *Location) ([]File, error)
	GetStatuses(location *Location) ([]github.RepoStatus, error)
}

type API struct {
	client *github.Client
	statsd statsdw.Interface
}

func New(token string, cacheDir string, statsd statsdw.Interface) Interface {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	oauthTransport := oauth2.NewClient(ctx, ts).Transport

	cache := diskcache.New(cacheDir)

	cacheTransport := &httpcache.Transport{
		Transport: oauthTransport,
		Cache:     cache,

		MarkCachedResponses: true,
	}

	client := github.NewClient(&http.Client{
		Transport: cacheTransport,
	})

	return &API{
		client: client,
		statsd: statsd,
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
		SHA:      *ghBranch.Commit.SHA,
		Message:  *ghBranch.Commit.Commit.Message,
		Author:   *ghBranch.Commit.Commit.Author.Name,
		Date:     *ghBranch.Commit.Commit.Author.Date,
	}

	log.WithFields(log.Fields{
		"Owner":         location.Owner,
		"Repo":          location.Repo,
		"Ref":           location.Ref,
		"RateLimit":     resp.Rate.Limit,
		"RateRemaining": resp.Rate.Remaining,
		"RateReset":     resp.Rate.Reset,
		"FromCache":     resp.Header.Get(httpcache.XFromCache),
		"Branch":        *ghBranch.Name,
		"Author":        *ghBranch.Commit.Commit.Author.Name,
		"SHA":           *ghBranch.Commit.SHA,
		"Date":          ghBranch.Commit.Commit.Author.Date,
	}).Debug("fetched branch information")

	gh.statsd.Gauge("github.rate.remaining", resp.Rate.Remaining)

	return branch, err
}

type FileFuture struct {
	file chan File
	err  chan error
}

func (ff *FileFuture) Get() (File, error) {
	select {
	case file := <-ff.file:
		close(ff.file)
		close(ff.err)
		return file, nil
	case err := <-ff.err:
		close(ff.file)
		close(ff.err)
		return File{}, err
	}
}

func (gh *API) GetFileAsync(location *Location) *FileFuture {
	ff := &FileFuture{
		file: make(chan File, 1),
		err:  make(chan error, 1),
	}

	go func() {
		file, err := gh.GetFile(location)
		if err != nil {
			ff.err <- err
		}
		ff.file <- file
	}()

	return ff
}

type ErrNotAFile struct {
	Location *Location
}

func (err ErrNotAFile) Error() string {
	return fmt.Sprintf(
		"unable to fetch file '%v' from GitHub; probably it's a directory",
		err.Location)

}

func (gh *API) GetFile(location *Location) (File, error) {
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
		return File{}, errors.Wrapf(err,
			"unable to fetch file '%v' from GitHub", location)
	}

	if file == nil {
		return File{}, ErrNotAFile{Location: location}
	}

	content, err := file.GetContent()
	if err != nil {
		return File{}, errors.Wrapf(err,
			"unable to decode file '%v'", location)
	}

	log.WithFields(log.Fields{
		"Size":          *file.Size,
		"URL":           *file.HTMLURL,
		"RateLimit":     resp.Rate.Limit,
		"RateRemaining": resp.Rate.Remaining,
		"RateReset":     resp.Rate.Reset,
		"FromCache":     resp.Header.Get(httpcache.XFromCache),
	}).Debugf("found file")

	gh.statsd.Gauge("github.rate.remaining", resp.Rate.Remaining)

	return File{location.Path, content}, nil
}

func (gh *API) GetFiles(location *Location) ([]File, error) {
	log.WithFields(
		log.Fields(structs.Map(location)),
	).Debug("downloading directory from GitHub")

	// Need to remove the trailing slash, because otherwise GitHub answers with
	// a redirect which doesn't contain the ref param anymore.
	cleanPath := strings.TrimRight(location.Path, "/")

	_, dir, resp, err := gh.client.Repositories.GetContents(
		context.Background(),
		location.Owner, location.Repo, cleanPath,
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
		"FromCache":     resp.Header.Get(httpcache.XFromCache),
	}).Debug("found files in directory")

	gh.statsd.Gauge("github.rate.remaining", resp.Rate.Remaining)

	futures := make(map[string]*FileFuture)
	for _, file := range dir {
		futures[*file.Name] = gh.GetFileAsync(&Location{
			Owner: location.Owner,
			Repo:  location.Repo,
			Path:  path.Join(location.Path, *file.Name),
			Ref:   location.Ref,
		})
	}

	var files []File
	for _, future := range futures {
		file, err := future.Get()
		if err != nil {
			_, notAFile := err.(ErrNotAFile)
			if !notAFile {
				return nil, errors.Wrapf(err,
					"unable to decode file '%v'",
					location)
			}

		}
		files = append(files, file)
	}

	return files, nil
}

func (gh *API) GetStatuses(location *Location) ([]github.RepoStatus, error) {
	log.WithFields(
		log.Fields(structs.Map(location)),
	).Debug("getting build status from GitHub")

	combined, resp, err := gh.client.Repositories.GetCombinedStatus(
		context.Background(),
		location.Owner, location.Repo, location.Ref,
		nil,
	)

	if err != nil {
		return nil, errors.Wrapf(err,
			"unable to fetch status '%v' from GitHub", location)
	}

	log.WithFields(log.Fields{
		"TotalStatuses": combined.GetTotalCount(),
		"SHA":           combined.GetSHA(),
		"RateLimit":     resp.Rate.Limit,
		"RateRemaining": resp.Rate.Remaining,
		"RateReset":     resp.Rate.Reset,
		"FromCache":     resp.Header.Get(httpcache.XFromCache),
	}).Debug("retrieved statuses page")
	gh.statsd.Gauge("github.rate.remaining", resp.Rate.Remaining)

	return combined.Statuses, nil
}
