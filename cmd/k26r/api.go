package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	"k8s.io/client-go/kubernetes"
)

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
