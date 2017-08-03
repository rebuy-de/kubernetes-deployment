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
	case *v1beta1extensions.Deployment:
		i.AddToDeployment(typed)
		return typed, nil

	case *v1beta1apps.StatefulSet:
		i.AddToStatefulSet(typed)
		return typed, nil

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support adding preStop hook")
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

func (i *Interceptor) AddToPodTemplace(tpl *v1.PodTemplateSpec) {
	for j := range tpl.Spec.Containers {
		i.AddToContainer(&tpl.Spec.Containers[j])
	}
}

func (i *Interceptor) AddToDeployment(deployment *v1beta1extensions.Deployment) {
	i.AddToPodTemplace(&deployment.Spec.Template)
}

func (i *Interceptor) AddToStatefulSet(set *v1beta1apps.StatefulSet) {
	i.AddToPodTemplace(&set.Spec.Template)
}
