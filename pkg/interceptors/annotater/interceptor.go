package annotater

import (
	"fmt"

	"github.com/benbjohnson/clock"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

type Interceptor struct {
	clock  clock.Clock
	branch gh.Branch
}

func New() *Interceptor {
	return &Interceptor{
		clock: clock.New(),
	}
}

func (i *Interceptor) PostFetch(branch *gh.Branch) error {
	i.branch = *branch
	return nil
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	now := i.clock.Now()

	annotations := accessor.GetAnnotations()
	if annotations == nil {
		annotations = map[string]string{}
	}

	annotations["rebuy.com/kubernetes-deployment.deployment-date"] = fmt.Sprint(now)
	annotations["rebuy.com/kubernetes-deployment.commit-sha"] = i.branch.SHA
	annotations["rebuy.com/kubernetes-deployment.commit-date"] = fmt.Sprint(i.branch.Date)
	annotations["rebuy.com/kubernetes-deployment.commit-author"] = i.branch.Author
	annotations["rebuy.com/kubernetes-deployment.commit-message"] = i.branch.Message
	annotations["rebuy.com/kubernetes-deployment.commit-location"] = i.branch.Location.String()

	accessor.SetAnnotations(annotations)

	return obj, nil
}
