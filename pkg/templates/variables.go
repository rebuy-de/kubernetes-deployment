package templates

import "github.com/imdario/mergo"

type Variables map[string]string

func (v *Variables) Defaults(defaults Variables) {
	mergo.Merge(v, defaults)
}
