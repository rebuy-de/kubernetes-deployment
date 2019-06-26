package rmresspec

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
)

type Interceptor struct {
}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	template := kubeutil.PodTemplateSpecAccessor(obj)
	if template != nil {
		RemoveFromPodTemplate(template)
	}

	return obj, nil
}

func RemoveFromPodTemplate(tpl *v1.PodTemplateSpec) {
	for i := range tpl.Spec.Containers {
		tpl.Spec.Containers[i].Resources = v1.ResourceRequirements{}
	}

	for i := range tpl.Spec.InitContainers {
		tpl.Spec.InitContainers[i].Resources = v1.ResourceRequirements{}
	}
}
