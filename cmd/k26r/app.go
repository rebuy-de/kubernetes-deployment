package main

import (
	"context"
	"net/http"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type App struct {
	HTTPListenAddress string
	GitHubToken       string
	SettingsFile      string
	Kubeconfig        string
}

func (app *App) Bind(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&app.HTTPListenAddress, "listen-addr", ":8080",
		"listening address for the HTTP server")

	cmd.PersistentFlags().StringVar(&app.SettingsFile, "settings-file", "",
		"path to the file containing the project settings")

	cmd.PersistentFlags().StringVar(&app.Kubeconfig, "kubeconfig", "",
		"path to the kubeconfig file to use for deploying")

	cmd.PersistentFlags().StringVar(&app.GitHubToken, "github-token", "",
		"authentication token for the GitHub API")
}

func (app *App) Run(ctx context.Context, cmd *cobra.Command, args []string) {
	var err error

	api := new(API)
	api.GitHub = gh.New(app.GitHubToken, statsdw.NullClient{})

	api.Kubectl = kubectl.New("kubectl", app.Kubeconfig)

	api.Kubernetes, err = newKubernetesClient(app.Kubeconfig)
	cmdutil.Must(err)

	api.Settings, err = settings.Read(api.Kubernetes)
	cmdutil.Must(err)

	server := &http.Server{
		Addr:    app.HTTPListenAddress,
		Handler: api.Mux(),
	}

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		logrus.Error(server.ListenAndServe())
		cancel()
	}()
	logrus.Info("started HTTP server")

	<-ctx.Done()
	logrus.Warn("shutting down HTTP server")
	logrus.Error(server.Shutdown(context.Background()))
}
