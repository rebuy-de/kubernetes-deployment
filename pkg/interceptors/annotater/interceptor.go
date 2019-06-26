package annotater

import (
	"fmt"
	"strings"
	"time"

	"github.com/benbjohnson/clock"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
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

	subObjects := kubeutil.SubObjectAccessor(obj)
	for j, sub := range subObjects {
		i.annotate(workload.GetName(), sub, j == 0)
	}

	return obj, nil
}

func (i *Interceptor) annotate(workload string, obj v1meta.Object, isRoot bool) {
	key := func(n string) string { return fmt.Sprintf("rebuy.com/kubernetes-deployment.%s", n) }

	template, ok := obj.(*core_v1.PodTemplateSpec)
	if ok {
		fmt.Println(template.Spec.Containers)
	}

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
