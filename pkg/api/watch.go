package api

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/pkg/api"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
)

const (
	ErrImagePull = "ErrImagePull"
	Error        = "Error"
)

func (app *App) Watch(selector string) error {
	log.WithFields(log.Fields{
		"LabelSelector": selector,
	}).Debugf("watching project")

	deployments, err := app.Clients.Kubernetes.
		Extensions().
		Deployments("").
		List(metav1.ListOptions{
			LabelSelector: selector,
		})
	if err != nil {
		return errors.Wrapf(err, "unable to list deployments")
	}

	for _, deployment := range deployments.Items {
		log.WithFields(log.Fields{
			"Name": deployment.ObjectMeta.Name,
		}).Debugf("watching deployment")

		ctx, done := context.WithCancel(context.Background())

		go func() {
			podc, err := app.WatchDeploymentPods(ctx, &deployment)
			if err != nil {
				log.Error(err)
				return
			}

			for pod := range podc {
				log.WithFields(log.Fields{
					"Name":      pod.ObjectMeta.Name,
					"Namespace": pod.ObjectMeta.Namespace,
				}).Debugf("pod changed")
				err := podReady(pod)
				if err != nil {
					log.WithFields(log.Fields{
						"Name":      pod.ObjectMeta.Name,
						"Namespace": pod.ObjectMeta.Namespace,
						"Status":    pod.Status,
					}).Warn(err)
				}
			}
		}()

		for deployment := range app.WatchDeployment(ctx, deployment.ObjectMeta.Name) {
			if deploymentRolloutComplete(deployment) {
				done()
			}
		}

		done()

		log.WithFields(log.Fields{
			"Name": deployment.ObjectMeta.Name,
		}).Debugf("deployment succeeded")
	}

	log.WithFields(log.Fields{
		"LabelSelector": selector,
	}).Debugf("deployment succeeded")

	return nil
}

func (app *App) GetReplicaSetForDeployment(deployment *v1beta1.Deployment) (*v1beta1.ReplicaSet, error) {
	replicaSets, err := app.Clients.Kubernetes.
		Extensions().
		ReplicaSets(deployment.ObjectMeta.Namespace).
		List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "unable to list replica sets")
	}

	deploymentRevision, ok := deployment.ObjectMeta.Annotations["deployment.kubernetes.io/revision"]
	if !ok {
		return nil, errors.Errorf("deployment doesn't have a revision annotation")
	}

	for _, rs := range replicaSets.Items {
		rsRevision, ok := rs.ObjectMeta.Annotations["deployment.kubernetes.io/revision"]
		if !ok {
			continue
		}

		if deploymentRevision != rsRevision {
			continue
		}

		for _, or := range rs.ObjectMeta.OwnerReferences {
			if or.UID == deployment.UID {
				return &rs, nil
			}
		}
	}

	return nil, errors.Errorf("could not found replicaset for deployment")
}

func (app *App) WatchDeploymentPods(ctx context.Context, deployment *v1beta1.Deployment) (chan *v1.Pod, error) {
	rs, err := app.GetReplicaSetForDeployment(deployment)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	log.WithFields(log.Fields{
		"Name":      rs.ObjectMeta.Name,
		"Namespace": rs.ObjectMeta.Namespace,
	}).Debug("found active replica set")

	lw := cache.NewListWatchFromClient(
		app.Clients.Kubernetes.Core().RESTClient(),
		"pods",
		api.NamespaceAll,
		fields.Everything())

	stop := make(chan struct{}, 1)
	pods := make(chan *v1.Pod)
	results := make(chan *v1.Pod)

	store, controller := cache.NewInformer(
		lw,
		&v1.Pod{},
		60*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				pods <- obj.(*v1.Pod)
			},
			UpdateFunc: func(old, obj interface{}) {
				pods <- obj.(*v1.Pod)
			},
			DeleteFunc: func(obj interface{}) {
				pods <- obj.(*v1.Pod)
			},
		})

	for _, obj := range store.List() {
		pods <- obj.(*v1.Pod)
	}

	go controller.Run(stop)

	go func() {
		for pod := range pods {
			for _, or := range pod.ObjectMeta.OwnerReferences {
				if or.UID == rs.UID {
					results <- pod
					break
				}
			}
		}
	}()

	go func() {
		<-ctx.Done()
		close(stop)
		close(pods)
		close(results)
	}()

	return results, nil
}

func (app *App) WatchDeployment(ctx context.Context, name string) chan *v1beta1.Deployment {
	lw := cache.NewListWatchFromClient(
		app.Clients.Kubernetes.ExtensionsV1beta1().RESTClient(),
		"deployments",
		api.NamespaceAll,
		fields.OneTermEqualSelector("metadata.name", name))

	stop := make(chan struct{}, 1)
	results := make(chan *v1beta1.Deployment)

	store, controller := cache.NewInformer(
		lw,
		&v1beta1.Deployment{},
		60*time.Second,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				results <- obj.(*v1beta1.Deployment)
			},
			UpdateFunc: func(old, obj interface{}) {
				results <- obj.(*v1beta1.Deployment)
			},
			DeleteFunc: func(obj interface{}) {
				results <- obj.(*v1beta1.Deployment)
			},
		})

	for _, obj := range store.List() {
		results <- obj.(*v1beta1.Deployment)
	}

	go controller.Run(stop)

	go func() {
		<-ctx.Done()
		close(stop)
		close(results)
	}()

	return results
}

func deploymentRolloutComplete(deployment *v1beta1.Deployment) bool {
	logger := log.WithFields(log.Fields{
		"Namespace":          deployment.ObjectMeta.Namespace,
		"Name":               deployment.ObjectMeta.Name,
		"ResourceVersion":    deployment.ObjectMeta.ResourceVersion,
		"UpdatedReplicas":    deployment.Status.UpdatedReplicas,
		"DesiredReplicas":    *(deployment.Spec.Replicas),
		"ActualGeneration":   deployment.ObjectMeta.Generation,
		"ObservedGeneration": deployment.Status.ObservedGeneration,
	})

	if deployment.Status.UpdatedReplicas == *(deployment.Spec.Replicas) &&
		deployment.Status.Replicas == *(deployment.Spec.Replicas) &&
		deployment.Status.AvailableReplicas == *(deployment.Spec.Replicas) &&
		deployment.Status.ObservedGeneration >= deployment.Generation {
		logger.Debug("deployment is up to date")
		return true
	}

	logger.Debug("rollout still in progress")
	return false
}

func podReady(pod *v1.Pod) error {
	for _, cs := range pod.Status.InitContainerStatuses {
		err := containerReady(cs)
		if err != nil {
			return err
		}
	}

	for _, cs := range pod.Status.ContainerStatuses {
		err := containerReady(cs)
		if err != nil {
			return err
		}
	}

	return nil
}

func containerReady(status v1.ContainerStatus) error {
	if status.State.Waiting != nil && status.State.Waiting.Reason == ErrImagePull {
		return fmt.Errorf("failed to pull docker image")
	}

	if status.State.Terminated != nil && status.State.Terminated.Reason == Error {
		return fmt.Errorf(
			"failed to start container (Container: %v, ExitCode: %v, Restarts: %v)",
			status.Name, status.State.Terminated.ExitCode, status.RestartCount)
	}

	return nil
}
