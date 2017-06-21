package api

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
)

type Clients struct {
	GitHub     gh.Client
	Kubernetes kubectl.Interface
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
	app.Clients = &Clients{
		GitHub:     gh.New(p.GitHubToken, p.HTTPCacheDir),
		Kubernetes: kubectl.New(p.KubectlPath, p.Kubeconfig),
	}

	app.Settings, err = settings.Read(p.Filename, app.Clients.GitHub)
	if err != nil {
		return nil, err
	}

	return app, nil
}
