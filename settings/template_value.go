package settings

type TemplateValue map[string]string

func (tv TemplateValue) Merge(other TemplateValue) TemplateValue {
	result := make(TemplateValue)

	for k, v := range tv {
		result[k] = v
	}

	for k, v := range other {
		result[k] = v
	}

	return result
}
