package gh

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

type Location struct {
	Owner, Repo, Path, Ref string `yaml:",omitempty"`
}

func NewLocation(location string) (*Location, error) {
	matches := ContentLocationRE.FindStringSubmatch(location)
	if matches == nil {
		return nil, errors.Errorf(
			"GitHub location must have the form `github.com/:owner:/:repo:/:path:`")
	}

	return &Location{
		Owner: matches[1],
		Repo:  matches[2],
		Path:  matches[3],
		Ref:   "master",
	}, nil
}

func (l Location) String() string {
	return fmt.Sprintf("github.com/%s/%s/%s@%s", l.Owner, l.Repo, l.Path, l.Ref)
}

func (l *Location) Defaults(defaults Location) {
	mergo.Merge(l, defaults)
}
