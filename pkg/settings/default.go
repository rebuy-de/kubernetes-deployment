package settings

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/imagechecker"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/prestopsleep"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/statuschecker"
)

var (
	Defaults = Service{
		Aliases: []string{},
		Location: gh.Location{
			Ref: "master",
		},
		Interceptors: Interceptors{
			PreStopSleep: PreStopSleepInterceptor{
				Options: prestopsleep.Options{
					Seconds: 3,
				},
			},
			GHStatusChecker: GHStatusCheckerInterceptor{
				Options: statuschecker.DefaultOptions,
			},
			ImageChecker: ImageCheckerInterceptor{
				Options: imagechecker.DefaultOptions,
			},
		},
	}
)
