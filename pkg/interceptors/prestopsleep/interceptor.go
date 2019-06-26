package prestopsleep

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
)

type Interceptor struct {
	SleepSeconds int
}

func New(sleepSeconds int) *Interceptor {
	return &Interceptor{
		SleepSeconds: sleepSeconds,
	}
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	template := kubeutil.PodTemplateSpecAccessor(obj)
	if template != nil {
		i.AddToPodTemplate(template)
	}

	return obj, nil
}

func (i *Interceptor) AddToContainer(container *v1.Container) {
	if container.Lifecycle == nil {
		container.Lifecycle = &v1.Lifecycle{}
	}

	if container.Lifecycle.PreStop != nil {
		return
	}

	container.Lifecycle.PreStop = &v1.Handler{
		Exec: &v1.ExecAction{
			Command: []string{"sleep", fmt.Sprintf("%d", i.SleepSeconds)},
		},
	}
}

func (i *Interceptor) AddToPodTemplate(tpl *v1.PodTemplateSpec) {
	for j := range tpl.Spec.Containers {
		i.AddToContainer(&tpl.Spec.Containers[j])
	}
}
