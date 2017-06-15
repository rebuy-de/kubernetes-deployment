package settings

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
)

type Service struct {
	Name      string              `yaml:"name,omitempty"`
	Location  gh.Location         `yaml:",inline"`
	Variables templates.Variables `yaml:"variables"`
}
