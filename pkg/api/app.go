package api

import (
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/prestopsleep"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/rmoldjob"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/rmresspec"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/statuschecker"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/waiter"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Clients struct {
	GitHub     gh.Interface
	Kubectl    kubectl.Interface
	Kubernetes kubernetes.Interface
	Statsd     statsdw.Interface
}

type App struct {
	Parameters   *Parameters
	Clients      *Clients
	Settings     *settings.Settings
	Interceptors *interceptors.Multi
}

func New(p *Parameters) (*App, error) {
	var err error

	app := new(App)

	app.Parameters = p

	app.Clients = &Clients{}
	app.Clients.Statsd = statsdw.New(p.StatsdAddress)
	app.Clients.GitHub = gh.New(p.GitHubToken, p.HTTPCacheDir, app.Clients.Statsd)
	app.Clients.Kubectl = kubectl.New(p.KubectlPath, p.Kubeconfig)

	app.Clients.Kubernetes, err = newKubernetesClient(p.Kubeconfig)
	if err != nil {
		return nil, err
	}

	app.Settings, err = settings.Read(p.Filename, app.Clients.GitHub)
	if err != nil {
		return nil, err
	}

	app.Settings.Clean(p.Context)

	return app, nil
}

func (app *App) StartInterceptors(service *settings.Service) {
	app.Interceptors = interceptors.New()

	interceptors := service.Interceptors

	if interceptors.Waiter.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "waiter",
		}).Debug("enabling waiter interceptor")
		app.Interceptors.Add(waiter.NewDeploymentWaitInterceptor(app.Clients.Kubernetes))
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

	if interceptors.GHStatusChecker.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "ghStatusChecker",
			"Options":     interceptors.GHStatusChecker.Options,
		}).Debug("enabling ghStatusChecker interceptor")
		app.Interceptors.Add(statuschecker.New(
			app.Clients.GitHub,
			interceptors.GHStatusChecker.Options,
		))
	}

	if interceptors.RemoveOldJob.Enabled == settings.Enabled {
		log.WithFields(log.Fields{
			"Interceptor": "removeOldJob",
		}).Debug("enabling removeOldJob interceptor")
		app.Interceptors.Add(rmoldjob.New(
			app.Clients.Kubernetes,
		))
	}
}

func (app *App) CurrentContext() settings.Service {
	contextName := app.Parameters.Context
	if contextName == "" {
		contextName = app.Settings.Defaults.Context
		log.WithFields(log.Fields{
			"Context": contextName,
		}).Debug("no context set; using default")
	}

	context, ok := app.Settings.Contexts[contextName]
	if !ok {
		context = app.Settings.Defaults
		log.WithFields(log.Fields{
			"Context": contextName,
		}).Debug("context not found; falling back to defaults")
	}

	return context
}

func (app *App) Close() error {
	return app.Interceptors.Close()
}

func newKubernetesClient(kubeconfig string) (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load kubernetes config")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize kubernetes client")
	}

	return client, nil
}
