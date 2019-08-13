package settings

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/imagechecker"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/injector"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/prestopsleep"
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
			ImageChecker: ImageCheckerInterceptor{
				Options: imagechecker.DefaultOptions,
			},
			Injector: InjectorInterceptor{
				Options: injector.DefaultOptions,
			},
		},
	}
)
