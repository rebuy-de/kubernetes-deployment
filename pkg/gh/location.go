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
	ContentLocationRE = regexp.MustCompile(`^github.com/([^/]+)/([^/]+)(/.*)?(@([^@]+))??$`)
)

type Location struct {
	Owner, Repo, Path, Ref string `yaml:",omitempty"`
}

func NewLocation(location string) (*Location, error) {
	formatError := errors.Errorf(
		"GitHub location must have the form `github.com/:owner:/:repo:/:path:[@:ref:]`, but got `%s`", location)

	result := new(Location)

	pathAndRef := strings.SplitN(location, "@", 2)
	if len(pathAndRef) != 1 && len(pathAndRef) != 2 {
		return nil, formatError
	}

	path := pathAndRef[0]
	if len(pathAndRef) > 1 {
		result.Ref = pathAndRef[1]
	}

	parts := strings.SplitN(path, "/", 4)
	if parts[0] != "github.com" || len(parts) < 3 {
		return nil, errors.Errorf(
			"GitHub location must have the form `github.com/:owner:/:repo:/:path:[@:ref:]`, but got `%s`", location)

	}

	result.Owner = parts[1]
	result.Repo = parts[2]

	// avoid index out of range error when skipping path
	parts = append(parts, "")
	result.Path = parts[3]

	return result, nil
}

func (l Location) String() string {
	path := fmt.Sprintf("github.com/%s/%s/%s", l.Owner, l.Repo, l.Path)
	if l.Ref != "" {
		path = fmt.Sprintf("%s@%s", path, l.Ref)
	}
	return path
}

func (l Location) ShortString() string {
	return fmt.Sprintf("%s/%s/%s", l.Owner, l.Repo, l.Path)
}

func (l *Location) Defaults(defaults Location) {
	mergo.Merge(l, defaults)
}

func (l *Location) Clean() {
	l.Path = filepath.Clean(strings.Trim(l.Path, "/")) + "/"
}
