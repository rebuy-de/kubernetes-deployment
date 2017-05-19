package settings

import "github.com/rebuy-de/kubernetes-deployment/pkg/gh"

type Service struct {
	Name           string         `yaml:"name,omitempty"`
	Location       gh.Location    `yaml:",inline"`
	TemplateValues TemplateValues `yaml:"template-values"`
}
