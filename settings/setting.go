package settings

import "time"

type Settings struct {
	Kubeconfig           *string        `yaml:"kubeconfig"`
	Output               *string        `yaml:"output"`
	Sleep                *time.Duration `yaml:"sleep"`
	SkipShuffle          *bool          `yaml:"skip-shuffle"`
	RetrySleep           *time.Duration `yaml:"retry-sleep"`
	RetryCount           *int           `yaml:"retry-count"`
	IgnoreDeployFailures *bool          `yaml:"ignore-deploy-failures"`
	TemplateValues       TemplateValue  `yaml:"template-values"`
}
