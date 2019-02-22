package annotater

import (
	"fmt"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/sirupsen/logrus"
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
		clock: clock.New(),
	}
}

func (i *Interceptor) PostFetch(branch *gh.Branch) error {
	i.branch = *branch
	i.now = i.clock.Now().In(i.timezone)
	return nil
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	i.annotate(accessor, accessor)

	return obj, nil
}

func (i *Interceptor) annotate(owner, obj v1meta.Object) {
	key := func(n string) string { return fmt.Sprintf("rebuy.com/kubernetes-deployment.%s", n) }

	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations[key("deployment-date")] = i.now.Format(time.RFC3339Nano)
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

	name := owner.GetName()
	labelName, ok := labels["name"]
	if ok && name != labelName {
		logrus.Warnf("Existing label name '%s' not match manifest name '%s'",
			labelName, name)
	} else if !ok {
		labels["name"] = name
	}

	obj.SetLabels(labels)
}
