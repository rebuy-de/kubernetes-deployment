package settings

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
