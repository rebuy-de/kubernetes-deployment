package api

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
)

type Clients struct {
	GitHub  gh.Interface
	Kubectl kubectl.Interface
	Statsd  statsdw.Interface
}

type App struct {
	Parameters *Parameters
	Clients    *Clients
	Settings   *settings.Settings
}

func New(p *Parameters) (*App, error) {
	var err error

	app := new(App)

	app.Parameters = p

	app.Clients = &Clients{}
	app.Clients.Statsd, err = statsdw.New(p.StatsdAddress)
	if err != nil {
		return nil, err
	}
	app.Clients.GitHub = gh.New(p.GitHubToken, p.HTTPCacheDir, app.Clients.Statsd)
	app.Clients.Kubectl = kubectl.New(p.KubectlPath, p.Kubeconfig)

	app.Settings, err = settings.Read(p.Filename, app.Clients.GitHub)
	if err != nil {
		return nil, err
	}

	return app, nil
}
