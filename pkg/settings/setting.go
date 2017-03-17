package settings

import "time"

type Settings struct {
	Kubeconfig           string            `yaml:"kubeconfig"`
	Output               string            `yaml:"output"`
	Sleep                time.Duration     `yaml:"sleep"`
	SkipShuffle          bool              `yaml:"skip-shuffle" mapstructure:"skip-shuffle"`
	RetrySleep           time.Duration     `yaml:"retry-sleep" mapstructure:"retry-sleep"`
	RetryCount           int               `yaml:"retry-count" mapstructure:"retry-count"`
	IgnoreDeployFailures bool              `yaml:"ignore-deploy-failures" mapstructure:"ignore-deploy-failures"`
	TemplateValues       map[string]string `yaml:"template-values" mapstructure:"template-values"`
}
