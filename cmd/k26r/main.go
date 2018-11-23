package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	defer cmdutil.HandleExit()
	if err := NewRootCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func NewRootCommand() *cobra.Command {
	cmd := cmdutil.NewRootCommand(new(App))
	cmd.Short = "k26r (kubernetes-deployment-server)"
	return cmd
}

type App struct {
	HTTPListenAddress string
	GitHubToken       string
	SettingsFile      string
	HTTPCacheDir      string
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
	cmd.PersistentFlags().StringVar(&app.HTTPCacheDir, "http-cache-dir",
		"/tmp/kubernetes-deployment-cache",
		"cache directory for HTTP client requests")
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

func (app *App) MustReadSettings() *settings.Settings {
	settings, err := settings.ReadFromFile(app.SettingsFile)
	cmdutil.Must(err)
	settings.Clean("")
	return settings
}

func (app *App) Run(ctx context.Context, cmd *cobra.Command, args []string) {
	var err error

	api := new(API)
	api.Settings = app.MustReadSettings()
	api.GitHub = gh.New(app.GitHubToken, app.HTTPCacheDir, statsdw.NullClient{})

	api.Kubectl = kubectl.New("kubectl", app.Kubeconfig)

	api.Kubernetes, err = newKubernetesClient(app.Kubeconfig)
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

type API struct {
	Settings   *settings.Settings
	GitHub     gh.Interface
	Kubectl    kubectl.Interface
	Kubernetes kubernetes.Interface
}

func (a *API) Mux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/dump/settings", a.HandleDumpSettings).
		Methods("GET")

	r.HandleFunc("/render/{project:[a-z0-9-_/]+}", a.HandleRender).
		Methods("GET")

	r.HandleFunc("/deploy/{project:[a-z0-9-_/]+}", a.HandleDeploy).
		Methods("POST")

	return r
}

func (a *API) HandleDumpSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	err := e.Encode(a.Settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) HandleRender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]

	branch := r.FormValue("branch")
	if branch == "" {
		branch = "master"
	}

	legacy := api.App{
		Settings: a.Settings,
		Clients: &api.Clients{
			GitHub: a.GitHub,
			Statsd: statsdw.NullClient{},
		},
	}

	objs, err := legacy.Generate(project, branch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	err = e.Encode(objs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *API) HandleDeploy(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]

	branch := r.FormValue("branch")
	if branch == "" {
		branch = "master"
	}

	legacy := api.App{
		Settings: a.Settings,
		Clients: &api.Clients{
			GitHub:     a.GitHub,
			Statsd:     statsdw.NullClient{},
			Kubectl:    a.Kubectl,
			Kubernetes: a.Kubernetes,
		},
	}

	err := legacy.Apply(project, branch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
}
