package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
)

type Parameters struct {
	Kubeconfig  string
	KubectlPath string `mapstructure:"kubectl-path"`

	GitHubToken  string `mapstructure:"github-token"`
	HTTPCacheDir string `mapstructure:"http-cache-dir"`

	Filename string

	ghClient gh.Client
}

func (p *Parameters) GitHubClient() gh.Client {
	if p.ghClient == nil {
		p.ghClient = gh.New(p.GitHubToken, p.HTTPCacheDir)
	}

	return p.ghClient
}

func (p *Parameters) Kubectl() kubectl.Interface {
	return kubectl.New(p.KubectlPath, p.Kubeconfig)
}

func (p *Parameters) LoadSettings() *settings.Settings {
	sett, err := settings.Read(p.Filename, p.GitHubClient())
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"ServiceCount": len(sett.Services),
		"Filename":     p.Filename,
	}).Debug("loaded service file")

	sett.Clean()

	return sett
}
