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
	"k8s.io/apimachinery/pkg/runtime"

	log "github.com/sirupsen/logrus"
)

// PodTemplateSpecAccessor extracts the PodTemplateSpec from any object that
// contains it. Similar to
// https://godoc.org/k8s.io/apimachinery/pkg/api/meta#Accessor
func PodTemplateSpecAccessor(obj runtime.Object) *core_v1.PodTemplateSpec {
	switch typed := obj.(type) {

	// StatefulSets
	case *apps_v1.StatefulSet:
		return &typed.Spec.Template
	case *apps_v1beta2.StatefulSet:
		return &typed.Spec.Template
	case *apps_v1beta1.StatefulSet:
		return &typed.Spec.Template

	// Deployments
	case *apps_v1.Deployment:
		return &typed.Spec.Template
	case *apps_v1beta2.Deployment:
		return &typed.Spec.Template
	case *apps_v1beta1.Deployment:
		return &typed.Spec.Template
	case *extensions_v1beta1.Deployment:
		return &typed.Spec.Template

	// DaemonSets
	case *apps_v1.DaemonSet:
		return &typed.Spec.Template
	case *apps_v1beta2.DaemonSet:
		return &typed.Spec.Template
	case *extensions_v1beta1.DaemonSet:
		return &typed.Spec.Template

	// CronJob
	case *batch_v1beta1.CronJob:
		return &typed.Spec.JobTemplate.Spec.Template
	case *batch_v2alpha1.CronJob:
		return &typed.Spec.JobTemplate.Spec.Template

	// Job
	case *batch_v1.Job:
		return &typed.Spec.Template

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support accessing of PodTemplateSpec")
		return nil
	}
}
