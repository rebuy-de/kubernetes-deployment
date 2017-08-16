package rmoldjob

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	v1batch "k8s.io/client-go/pkg/apis/batch/v1"
)

const Annotation = "rebuy.com/delete-on-deploy"

type Interceptor struct {
	client kubernetes.Interface
}

func New(client kubernetes.Interface) *Interceptor {
	return &Interceptor{
		client: client,
	}
}

func (i *Interceptor) PreManifestApply(obj runtime.Object) (runtime.Object, error) {
	job, ok := obj.(*v1batch.Job)
	if !ok {
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support deletion of old jobs")
		return obj, nil
	}

	name := job.ObjectMeta.Name
	namespace := job.ObjectMeta.Namespace

	if namespace == "" {
		namespace = "default"
	}

	has, err := i.hasJob(namespace, name)
	if !has || err != nil {
		return obj, errors.WithStack(err)
	}

	if !v1meta.HasAnnotation(job.ObjectMeta, Annotation) ||
		job.ObjectMeta.Annotations[Annotation] != "true" {
		log.WithFields(log.Fields{
			"Namespace": namespace,
			"Name":      name,
		}).Debug("skip deletion of job, because the annotation is missing")
		return obj, nil
	}

	log.WithFields(log.Fields{
		"Namespace": namespace,
		"Name":      name,
	}).Infof("deleting pending job: %s", name)

	return obj, i.client.
		Batch().
		Jobs(namespace).
		Delete(name, nil)
}

func (i *Interceptor) hasJob(namespace, name string) (bool, error) {
	_, err := i.client.
		Batch().
		Jobs(namespace).
		Get(name, v1meta.GetOptions{})

	if err == nil {
		return true, nil
	}

	status, ok := err.(*k8serr.StatusError)
	if !ok {
		return false, err
	}

	log.WithFields(log.Fields{
		"Error": fmt.Sprintf("%#v", status.ErrStatus),
	}).Debug("got error status")

	if status.ErrStatus.Reason == v1meta.StatusReasonNotFound {
		return false, nil
	}

	return false, err
}
