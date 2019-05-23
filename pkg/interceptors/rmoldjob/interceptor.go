package rmoldjob

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	v1batch "k8s.io/api/batch/v1"
	k8serr "k8s.io/apimachinery/pkg/api/errors"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

const (
	AnnotationDelete  = "rebuy.com/delete-on-deploy"
	AnnotationCascade = "rebuy.com/cascade-delete"
)

type Interceptor struct {
	client kubernetes.Interface
}

func New(client kubernetes.Interface) *Interceptor {
	return &Interceptor{
		client: client,
	}
}

func (i *Interceptor) PreManifestApply(obj runtime.Object) (runtime.Object, error) {
	logger := log.WithField("Type", reflect.TypeOf(obj))

	job, ok := obj.(*v1batch.Job)
	if !ok {
		logger.Debug("type doesn't support deletion of old jobs")
		return obj, nil
	}

	name := job.ObjectMeta.Name
	namespace := job.ObjectMeta.Namespace

	if namespace == "" {
		namespace = "default"
	}

	logger = logger.WithField("Namespace", namespace)
	logger = logger.WithField("Name", name)

	has, err := i.hasJob(namespace, name)
	if !has || err != nil {
		return obj, errors.WithStack(err)
	}

	if !v1meta.HasAnnotation(job.ObjectMeta, AnnotationDelete) ||
		job.ObjectMeta.Annotations[AnnotationDelete] != "true" {
		logger.Debug("skip deletion of job, because the annotation is missing")
		return obj, nil
	}

	logger.Infof("deleting pending job: %s", name)

	opts := new(v1meta.DeleteOptions)

	if v1meta.HasAnnotation(job.ObjectMeta, AnnotationCascade) &&
		job.ObjectMeta.Annotations[AnnotationCascade] == "true" {
		logger.Debug("enabling cascading delete")
		policy := v1meta.DeletePropagationBackground
		opts.PropagationPolicy = &policy
	}

	return obj, i.client.
		BatchV1().
		Jobs(namespace).
		Delete(name, opts)
}

func (i *Interceptor) hasJob(namespace, name string) (bool, error) {
	_, err := i.client.
		BatchV1().
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
