package settings

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/prestopsleep"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/statuschecker"
)

type Interceptors struct {
	PreStopSleep        PreStopSleepInterceptor    `yaml:"preStopSleep"`
	RemoveResourceSpecs Interceptor                `yaml:"removeResourceSpecs"`
	RemoveOldJob        Interceptor                `yaml:"removeOldJob"`
	Waiter              Interceptor                `yaml:"waiter"`
	GHStatusChecker     GHStatusCheckerInterceptor `yaml:"ghStatusChecker"`
}

type Interceptor struct {
	Enabled TriState `yaml:"enabled"`
}

type PreStopSleepInterceptor struct {
	Enabled TriState             `yaml:"enabled"`
	Options prestopsleep.Options `yaml:"options"`
}

type GHStatusCheckerInterceptor struct {
	Enabled TriState              `yaml:"enabled"`
	Options statuschecker.Options `yaml:"options"`
}
