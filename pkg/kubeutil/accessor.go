package kubeutil

import (
	"reflect"

	apps_v1 "k8s.io/api/apps/v1"
	apps_v1beta1 "k8s.io/api/apps/v1beta1"
	apps_v1beta2 "k8s.io/api/apps/v1beta2"
	batch_v1 "k8s.io/api/batch/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	batch_v2alpha1 "k8s.io/api/batch/v2alpha1"
	core_v1 "k8s.io/api/core/v1"
	extensions_v1beta1 "k8s.io/api/extensions/v1beta1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	log "github.com/sirupsen/logrus"
)

func SubObjectAccessor(obj runtime.Object) []v1meta.Object {
	metaObject, ok := obj.(v1meta.Object)
	if !ok {
		return nil
	}

	objs := []v1meta.Object{metaObject}

	switch typed := obj.(type) {

	// StatefulSets
	case *apps_v1.StatefulSet:
		objs = append(objs, &typed.Spec.Template)
	case *apps_v1beta2.StatefulSet:
		objs = append(objs, &typed.Spec.Template)
	case *apps_v1beta1.StatefulSet:
		objs = append(objs, &typed.Spec.Template)

	// Deployments
	case *apps_v1.Deployment:
		objs = append(objs, &typed.Spec.Template)
	case *apps_v1beta2.Deployment:
		objs = append(objs, &typed.Spec.Template)
	case *apps_v1beta1.Deployment:
		objs = append(objs, &typed.Spec.Template)
	case *extensions_v1beta1.Deployment:
		objs = append(objs, &typed.Spec.Template)

	// DaemonSets
	case *apps_v1.DaemonSet:
		objs = append(objs, &typed.Spec.Template)
	case *apps_v1beta2.DaemonSet:
		objs = append(objs, &typed.Spec.Template)
	case *extensions_v1beta1.DaemonSet:
		objs = append(objs, &typed.Spec.Template)

	// CronJob
	case *batch_v1beta1.CronJob:
		objs = append(objs, &typed.Spec.JobTemplate)
		objs = append(objs, &typed.Spec.JobTemplate.Spec.Template)
	case *batch_v2alpha1.CronJob:
		objs = append(objs, &typed.Spec.JobTemplate)
		objs = append(objs, &typed.Spec.JobTemplate.Spec.Template)

	// Job
	case *batch_v1.Job:
		objs = append(objs, &typed.Spec.Template)

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support accessing of sub objects")
	}

	return objs
}

// PodTemplateSpecAccessor extracts the PodTemplateSpec from any object that
// contains it. Similar to
// https://godoc.org/k8s.io/apimachinery/pkg/api/meta#Accessor
func PodTemplateSpecAccessor(obj runtime.Object) *core_v1.PodTemplateSpec {
	objects := SubObjectAccessor(obj)

	if objects == nil || len(objects) < 1 {
		return nil
	}

	// A PodTemplateSpec does not have any sub objects. Therefore it can only
	// be the last one in the list.
	sobj := objects[len(objects)-1]

	spec, ok := sobj.(*core_v1.PodTemplateSpec)
	if !ok {
		return nil
	}

	return spec
}
