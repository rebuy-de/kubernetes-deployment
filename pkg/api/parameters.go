package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
)

type Parameters struct {
	Kubeconfig  string
	GitHubToken string `mapstructure:"github-token"`
	Filename    string
}

func (p *Parameters) GitHubClient() gh.Client {
	return gh.New(p.GitHubToken)
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
