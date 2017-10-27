package settings

import (
	"strings"

	"github.com/imdario/mergo"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

type Service struct {
	Name         string              `yaml:"name,omitempty"`
	Aliases      []string            `yaml:"aliases,omitempty"`
	Context      string              `yaml:"context,omitempty"`
	Location     gh.Location         `yaml:",inline,omitempty"`
	Variables    templates.Variables `yaml:"variables,omitempty"`
	Interceptors Interceptors        `yaml:"interceptors,omitempty"`
}

func (s *Service) Defaults(defaults Service) {
	mergo.Merge(s, defaults)
}

func (s *Service) Clean(defaults Service) {
	s.Location.Clean()
	defaults.Location.Clean()

	if strings.TrimSpace(s.Name) == "" {
		nameParts := []string{}
		if s.Location.Owner != defaults.Location.Owner {
			nameParts = append(nameParts, s.Location.Owner)
		}

		if s.Location.Repo != defaults.Location.Repo {
			nameParts = append(nameParts, s.Location.Repo)
		}

		if s.Location.Path != defaults.Location.Path {
			path := s.Location.Path
			path = strings.Trim(path, "/")
			path = strings.Replace(path, "/", "-", -1)
			nameParts = append(nameParts, path)
		}

		s.Name = strings.Join(nameParts, "-")
	}
}
