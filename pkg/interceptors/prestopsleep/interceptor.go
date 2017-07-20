package prestopsleep

import (
	"fmt"
	"reflect"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/pkg/api/v1"
	v1beta1apps "k8s.io/client-go/pkg/apis/apps/v1beta1"
	v1beta1extensions "k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

const SleepSeconds = 5

type Interceptor struct {
}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) ManifestRendered(obj runtime.Object) (runtime.Object, error) {
	switch typed := obj.(type) {
	case *v1beta1extensions.Deployment:
		AddToDeployment(typed)
		return typed, nil

	case *v1beta1apps.StatefulSet:
		AddToStatefulSet(typed)
		return typed, nil

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support adding preStop hook")
	}
	return obj, nil
}

func AddToContainer(container *v1.Container) {
	if container.Lifecycle == nil {
		container.Lifecycle = &v1.Lifecycle{}
	}

	if container.Lifecycle.PreStop != nil {
		return
	}

	container.Lifecycle.PreStop = &v1.Handler{
		Exec: &v1.ExecAction{
			Command: []string{"sleep", fmt.Sprintf("%d", SleepSeconds)},
		},
	}
}

func AddToPodTemplace(tpl *v1.PodTemplateSpec) {
	for i := range tpl.Spec.Containers {
		AddToContainer(&tpl.Spec.Containers[i])
	}
}

func AddToDeployment(deployment *v1beta1extensions.Deployment) {
	AddToPodTemplace(&deployment.Spec.Template)
}

func AddToStatefulSet(set *v1beta1apps.StatefulSet) {
	AddToPodTemplace(&set.Spec.Template)
}
