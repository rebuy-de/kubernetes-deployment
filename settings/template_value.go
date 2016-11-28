package settings

type TemplateValues map[string]string

func (tv TemplateValues) Merge(other TemplateValues) TemplateValues {
	result := make(TemplateValues)

	for k, v := range tv {
		result[k] = v
	}

	for k, v := range other {
		result[k] = v
	}

	return result
}
