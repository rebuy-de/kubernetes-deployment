package settings

import "github.com/rebuy-de/kubernetes-deployment/pkg/gh"

var (
	Defaults = Service{
		Context: "default",
		Location: gh.Location{
			Ref: "master",
		},
		Interceptors: Interceptors{
			PreStopSleep: PreStopSleepInterceptor{
				Options: PreStopSleepOptions{
					Seconds: 3,
				},
			},
			GHStatusChecker: GHStatusCheckerInterceptor{
				Options: GHStatusCheckerOptions{
					TargetURLRegex: `.*`,
					JobRegex:       `.*`,
				},
			},
		},
	}
)
