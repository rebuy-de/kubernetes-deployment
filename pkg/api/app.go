package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/annotater"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/grafannotator"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/imagechecker"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/injector"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/prestopsleep"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/rmoldjob"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/rmresspec"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/waiter"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
)

type Clients struct {
}

type App struct {
	GitHub       gh.Interface
	AWS          *session.Session
	Kubectl      kubectl.Interface
	Kubernetes   kubernetes.Interface
	Statsd       statsdw.Interface
	Settings     *settings.Settings
	Interceptors *interceptors.Multi
}

func (app *App) StartInterceptors(service *settings.Service) {
	app.Interceptors = interceptors.New()

	interceptors := service.Interceptors

	if interceptors.Waiter.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "waiter",
		}).Debug("enabling waiter interceptor")
		app.Interceptors.Add(waiter.NewDeploymentWaitInterceptor(app.Kubernetes))
	}

	if interceptors.PreStopSleep.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor":  "preStopSleep",
			"SleepSeconds": interceptors.PreStopSleep.Options.Seconds,
		}).Debug("enabling preStopSleep interceptor")
		app.Interceptors.Add(prestopsleep.New(interceptors.PreStopSleep.Options.Seconds))
	}

	if interceptors.RemoveResourceSpecs.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "removeResourceSpecs",
		}).Debug("enabling removeResourceSpecs interceptor")
		app.Interceptors.Add(rmresspec.New())
	}

	if interceptors.ImageChecker.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "imageChecker",
			"Options":     interceptors.ImageChecker.Options,
		}).Debug("enabling imageChecker interceptor")
		app.Interceptors.Add(imagechecker.New(
			app.AWS,
			interceptors.ImageChecker.Options,
		))
	}

	if interceptors.RemoveOldJob.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "removeOldJob",
		}).Debug("enabling removeOldJob interceptor")
		app.Interceptors.Add(rmoldjob.New(
			app.Kubernetes,
		))
	}

	if interceptors.Annotater.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "annotater",
		}).Debug("enabling annotater interceptor")
		app.Interceptors.Add(annotater.New())
	}

	if interceptors.Grafannotator.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "grafannotator",
			"Options":     interceptors.Grafannotator.Options,
		}).Debug("enabling grafannotator interceptor")
		app.Interceptors.Add(grafannotator.New(
			interceptors.Grafannotator.Options,
		))
	}

	if interceptors.Injector.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "injector",
		}).Debug("enabling injector interceptor")
		app.Interceptors.Add(injector.New(
			interceptors.Injector.Options,
		))
	}
}

func (app *App) Close() error {
	return app.Interceptors.Close()
}
