package templates

import "github.com/imdario/mergo"

type Values map[string]string

func (v *Values) Defaults(defaults Values) {
	mergo.Merge(v, defaults)
}
