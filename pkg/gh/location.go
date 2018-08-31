package gh

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

var (
	ContentLocationRE = regexp.MustCompile(`^github.com/([^/]+)/([^/]+)/(.*?)(@([^@]+))?$`)
)

type Location struct {
	Owner, Repo, Path, Ref string `yaml:",omitempty"`
}

func NewLocation(location string) (*Location, error) {
	matches := ContentLocationRE.FindStringSubmatch(location)
	if matches == nil {
		return nil, errors.Errorf(
			"GitHub location must have the form `github.com/:owner:/:repo:/:path:[@:ref:]`, but got `%s`",
			location)
	}

	return &Location{
		Owner: matches[1],
		Repo:  matches[2],
		Path:  matches[3],
		Ref:   matches[5],
	}, nil
}

func (l Location) String() string {
	path := fmt.Sprintf("github.com/%s/%s/%s", l.Owner, l.Repo, l.Path)
	if l.Ref != "" {
		path = fmt.Sprintf("%s@%s", path, l.Ref)
	}
	return path
}

func (l *Location) Defaults(defaults Location) {
	mergo.Merge(l, defaults)
}

func (l *Location) Clean() {
	l.Path = filepath.Clean(strings.Trim(l.Path, "/")) + "/"
}
