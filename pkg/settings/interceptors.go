package settings

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/grafannotator"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/imagechecker"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/injector"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/prestopsleep"
)

type Interceptors struct {
	PreStopSleep        PreStopSleepInterceptor  `yaml:"preStopSleep"`
	RemoveResourceSpecs Interceptor              `yaml:"removeResourceSpecs"`
	RemoveOldJob        Interceptor              `yaml:"removeOldJob"`
	Waiter              Interceptor              `yaml:"waiter"`
	Annotater           Interceptor              `yaml:"annotater"`
	Grafannotator       GrafannotatorInterceptor `yaml:"grafannotator"`
	Injector            InjectorInterceptor      `yaml:"injector"`
	ImageChecker        ImageCheckerInterceptor  `yaml:"imageChecker"`
}

type Interceptor struct {
	Enabled TriState `yaml:"enabled"`
}

type PreStopSleepInterceptor struct {
	Enabled TriState             `yaml:"enabled"`
	Options prestopsleep.Options `yaml:"options"`
}

type InjectorInterceptor struct {
	Enabled TriState         `yaml:"enabled"`
	Options injector.Options `yaml:"options"`
}

type ImageCheckerInterceptor struct {
	Enabled TriState             `yaml:"enabled"`
	Options imagechecker.Options `yaml:"options"`
}

type GrafannotatorInterceptor struct {
	Enabled TriState              `yaml:"enabled"`
	Options grafannotator.Options `yaml:"options"`
}
