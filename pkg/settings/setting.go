package settings

import "time"

type Settings struct {
	Kubeconfig           string         `yaml:"kubeconfig"`
	Output               string         `yaml:"output"`
	Sleep                time.Duration  `yaml:"sleep"`
	SkipShuffle          bool           `yaml:"skip-shuffle" mapstructure:"skip-shuffle"`
	IgnoreDeployFailures bool           `yaml:"ignore-deploy-failures" mapstructure:"ignore-deploy-failures"`
	TemplateValues       TemplateValues `yaml:"template-values" mapstructure:"template-values"`
}

type TemplateValues []TemplateValue

func (tv TemplateValues) ToMap() map[string]string {
	result := make(map[string]string)

	for _, kv := range tv {
		result[kv.Name] = kv.Value
	}

	return result
}

type TemplateValue struct {
	Name  string `yaml:"name" mapstructure:"name"`
	Value string `yaml:"value" mapstructure:"value"`
}
