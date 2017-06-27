package api

import (
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors/waiter"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
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
	app.Clients.Statsd, err = statsdw.New(p.StatsdAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize statsd client")
	}

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

	app.Interceptors = interceptors.New(
		waiter.NewDeploymentWaitInterceptor(app.Clients.Kubernetes),
	)

	return app, nil
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
