package rmresspec

import (
	"reflect"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type Interceptor struct {
}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	switch typed := obj.(type) {
	case *v1beta1extensions.Deployment:
		return RemoveFromDeployment(typed), nil

	case *v1beta1apps.StatefulSet:
		return RemoveFromStatefulSet(typed), nil

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support removal of resource specs")
	}
	return obj, nil
}

func RemoveFromPodTemplace(tpl v1.PodTemplateSpec) v1.PodTemplateSpec {
	for i := range tpl.Spec.Containers {
		tpl.Spec.Containers[i].Resources = v1.ResourceRequirements{}
	}

	for i := range tpl.Spec.InitContainers {
		tpl.Spec.InitContainers[i].Resources = v1.ResourceRequirements{}
	}

	return tpl
}

func RemoveFromDeployment(deployment *v1beta1extensions.Deployment) *v1beta1extensions.Deployment {
	deployment.Spec.Template = RemoveFromPodTemplace(deployment.Spec.Template)
	return deployment
}

func RemoveFromStatefulSet(set *v1beta1apps.StatefulSet) *v1beta1apps.StatefulSet {
	set.Spec.Template = RemoveFromPodTemplace(set.Spec.Template)
	return set
}
