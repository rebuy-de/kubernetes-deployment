package annotater

import (
	"fmt"
	"strings"
	"time"

	"github.com/benbjohnson/clock"

	apps_v1 "k8s.io/api/apps/v1"
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	apps_v1beta2 "k8s.io/api/apps/v1beta2"
	batch_v1 "k8s.io/api/batch/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/meta"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

type Interceptor struct {
	clock    clock.Clock
	now      time.Time
	branch   gh.Branch
	timezone *time.Location
}

func New() *Interceptor {
	return &Interceptor{
		timezone: time.Local,
		clock:    clock.New(),
	}
}

func (i *Interceptor) PostFetch(branch *gh.Branch) error {
	i.branch = *branch
	i.now = i.clock.Now().In(i.timezone)
	return nil
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	workload, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	i.annotate(workload.GetName(), workload, true)

	switch typed := obj.(type) {
	case *apps_v1.Deployment:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *apps_v1beta2.Deployment:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *apps_v1beta1.Deployment:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *extensions_v1beta1.Deployment:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)

	case *apps_v1.DaemonSet:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *apps_v1beta2.DaemonSet:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *extensions_v1beta1.DaemonSet:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)

	case *apps_v1.StatefulSet:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *apps_v1beta2.StatefulSet:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	case *apps_v1beta1.StatefulSet:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)

	case *batch_v1beta1.CronJob:
		i.annotate(workload.GetName(), &typed.Spec.JobTemplate, false)
		i.annotate(workload.GetName(), &typed.Spec.JobTemplate.Spec.Template, false)

	case *batch_v1.Job:
		i.annotate(workload.GetName(), &typed.Spec.Template, false)
	}

	return obj, nil
}

func (i *Interceptor) annotate(workload string, obj v1meta.Object, isRoot bool) {
	key := func(n string) string { return fmt.Sprintf("rebuy.com/kubernetes-deployment.%s", n) }

	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	if isRoot {
		annotations[key("deployment-date")] = i.now.Format(time.RFC3339Nano)
	}
	annotations[key("commit-sha")] = i.branch.SHA
	annotations[key("commit-date")] = i.branch.Date.In(i.timezone).Format(time.RFC3339Nano)
	annotations[key("commit-author")] = i.branch.Author
	annotations[key("commit-message")] = i.branch.Message
	annotations[key("commit-location")] = i.branch.Location.String()

	obj.SetAnnotations(annotations)

	labels := obj.GetLabels()
	if labels == nil {
		labels = map[string]string{}
	}

	labels[key("workload-name")] = strings.Replace(workload, ":", "-", -1)
	labels[key("repo")] = i.branch.Location.Repo
	labels[key("branch")] = i.branch.Name

	obj.SetLabels(labels)
}
