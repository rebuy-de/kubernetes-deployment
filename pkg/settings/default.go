package settings

import "github.com/rebuy-de/kubernetes-deployment/pkg/gh"

var (
	Defaults = Service{
		Context: "default",
		Location: gh.Location{
			Ref: "master",
		},
	}
)
