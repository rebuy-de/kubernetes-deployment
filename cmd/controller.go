package cmd

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
)

type Controller struct {
	App *api.App
}

func (c *Controller) Mux() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/dump/settings", c.HandleDumpSettings).
		Methods("GET")

	r.HandleFunc("/render/{project:[a-z0-9-_/]+}", c.HandleRender).
		Methods("GET")

	return r
}

func (c *Controller) HandleDumpSettings(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	e := json.NewEncoder(w)
	e.SetIndent("", "    ")
	err := e.Encode(c.App.Settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Controller) HandleRender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	project := vars["project"]

	branch := r.FormValue("branch")
	if branch == "" {
		branch = "master"
	}

	objs, err := c.App.Generate(project, branch)
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
