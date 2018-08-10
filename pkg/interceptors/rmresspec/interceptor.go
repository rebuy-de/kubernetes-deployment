package rmresspec

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"

	v1apps "k8s.io/api/apps/v1"
	v1beta1apps "k8s.io/api/apps/v1beta1"
	v1beta2apps "k8s.io/api/apps/v1beta2"
	v1beta1extensions "k8s.io/api/extensions/v1beta1"
)

type Interceptor struct {
}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	switch typed := obj.(type) {
	case *v1apps.StatefulSet:
		RemoveFromPodTemplate(typed.Spec.Template)

	case *v1beta2apps.StatefulSet:
		RemoveFromPodTemplate(typed.Spec.Template)

	case *v1beta1apps.StatefulSet:
		RemoveFromPodTemplate(typed.Spec.Template)

	case *v1apps.Deployment:
		RemoveFromPodTemplate(typed.Spec.Template)

	case *v1beta2apps.Deployment:
		RemoveFromPodTemplate(typed.Spec.Template)

	case *v1beta1apps.Deployment:
		RemoveFromPodTemplate(typed.Spec.Template)

	case *v1beta1extensions.Deployment:
		RemoveFromPodTemplate(typed.Spec.Template)

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support removal of resource specs")
	}

	return obj, nil
}

func RemoveFromPodTemplate(tpl v1.PodTemplateSpec) {
	for i := range tpl.Spec.Containers {
		tpl.Spec.Containers[i].Resources = v1.ResourceRequirements{}
	}

	for i := range tpl.Spec.InitContainers {
		tpl.Spec.InitContainers[i].Resources = v1.ResourceRequirements{}
	}
}
