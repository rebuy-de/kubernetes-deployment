package settings

import (
	"github.com/imdario/mergo"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

type Service struct {
	Name                string              `yaml:"name,omitempty"`
	Context             string              `yaml:"context,omitempty"`
	Location            gh.Location         `yaml:",inline,omitempty"`
	Variables           templates.Variables `yaml:"variables,omitempty"`
	RemoveResourceSpecs bool                `yaml:"removeResourceSpecs,omitempty"` // deprecated
}

func (s *Service) Defaults(defaults Service) {
	mergo.Merge(s, defaults)
}
