package prestopsleep

import (
	"fmt"
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
	SleepSeconds int
}

func New(sleepSeconds int) *Interceptor {
	return &Interceptor{
		SleepSeconds: sleepSeconds,
	}
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	switch typed := obj.(type) {
	case *v1apps.StatefulSet:
		i.AddToPodTemplate(&typed.Spec.Template)

	case *v1beta2apps.StatefulSet:
		i.AddToPodTemplate(&typed.Spec.Template)

	case *v1beta1apps.StatefulSet:
		i.AddToPodTemplate(&typed.Spec.Template)

	case *v1apps.Deployment:
		i.AddToPodTemplate(&typed.Spec.Template)

	case *v1beta2apps.Deployment:
		i.AddToPodTemplate(&typed.Spec.Template)

	case *v1beta1apps.Deployment:
		i.AddToPodTemplate(&typed.Spec.Template)

	case *v1beta1extensions.Deployment:
		i.AddToPodTemplate(&typed.Spec.Template)

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support removal of resource specs")
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
