package settings

import "github.com/imdario/mergo"

type TemplateValues map[string]string

func (v *TemplateValues) Defaults(defaults TemplateValues) {
	mergo.Merge(v, defaults)
}
