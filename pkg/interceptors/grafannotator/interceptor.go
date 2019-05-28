package grafannotator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

type Query struct {
	Time int64    `json:"time"`
	Tags []string `json:"tags"`
	Text string   `json:"text"`
}

type Options struct {
	TargetURL string `yaml:"targetURL"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
}

type Interceptor struct {
	Options Options
}

func New(options Options) *Interceptor {
	return &Interceptor{
		Options: options,
	}
}

func (i *Interceptor) PostManifestApply(obj runtime.Object) error {
	var matchPR []string
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}

	kind := obj.GetObjectKind()
	if kind == nil {
		return nil
	}
	if !(kind.GroupVersionKind().Kind == "Deployment" ||
		kind.GroupVersionKind().Kind == "StatefulSet" ||
		kind.GroupVersionKind().Kind == "Daemonset" ||
		kind.GroupVersionKind().Kind == "CronJob" ||
		kind.GroupVersionKind().Kind == "Job") {
		return nil
	}

	labels := accessor.GetLabels()
	annotations := accessor.GetAnnotations()

	appname := labels["rebuy.com/kubernetes-deployment.workload-name"]
	repo := labels["rebuy.com/kubernetes-deployment.repo"]
	branch := labels["rebuy.com/kubernetes-deployment.branch"]
	commit := annotations["rebuy.com/kubernetes-deployment.commit-sha"]
	commitmessage := annotations["rebuy.com/kubernetes-deployment.commit-message"]
	date, err := time.Parse(time.RFC3339Nano, annotations["rebuy.com/kubernetes-deployment.deployment-date"])
	if err != nil {
		return err
	}

	text := "Branch: <a href='https://github.com/rebuy-de/" + repo + "/tree/" + branch + "'>" + branch + "</a></br>"
	if branch == "master" {
		regexPR := regexp.MustCompile(`Merge pull request #(\d+)`)
		matchPR = regexPR.FindStringSubmatch(commitmessage)
		if len(matchPR) > 1 {
			text += "Merged PR: <a href='https://github.com/rebuy-de/" + repo + "/pull/" + matchPR[1] + "'>#" + matchPR[1] + "</a></br>"
			regex := regexp.MustCompile(`.*\n.*\n(.*)`)
			match := regex.FindStringSubmatch(commitmessage)
			if len(match) > 1 {
				text += "Title: " + match[1]
			}
		}
	} else if branch != "master" || matchPR == nil {
		text += "Commit: <a href='https://github.com/rebuy-de/" + repo + "/commit/" + commit + "'>" + commitmessage + "</a></br>"
	}

	query := Query{
		Time: date.Unix() * 1000,
		Tags: []string{"deployment", "app=" + appname},
		Text: text,
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", i.Options.TargetURL+"/api/annotations", bytes.NewBuffer(queryJSON))
	if err != nil {
		return err
	}
	req.SetBasicAuth(i.Options.Username, i.Options.Password)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
