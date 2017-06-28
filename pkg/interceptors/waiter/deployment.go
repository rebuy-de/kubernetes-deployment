package waiter

import (
	"context"
	"fmt"
	"sync"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

type DeploymentWaitInterceptor struct {
	client kubernetes.Interface

	ctx       context.Context
	cancel    context.CancelFunc
	waitgroup *sync.WaitGroup
}

func NewDeploymentWaitInterceptor(client kubernetes.Interface) *DeploymentWaitInterceptor {
	ctx, cancel := context.WithCancel(context.Background())

	return &DeploymentWaitInterceptor{
		client:    client,
		ctx:       ctx,
		cancel:    cancel,
		waitgroup: new(sync.WaitGroup),
	}
}

func (dwi *DeploymentWaitInterceptor) AllManifestsApplied([]runtime.Object) error {
	dwi.waitgroup.Wait()
	return nil
}

func (dwi *DeploymentWaitInterceptor) Close() error {
	dwi.cancel()
	dwi.waitgroup.Wait()
	return nil
}

func (dwi *DeploymentWaitInterceptor) ManifestApplied(obj runtime.Object) error {
	deployment, ok := obj.(*v1beta1.Deployment)
	if !ok {
		return nil
	}

	log.WithFields(log.Fields{
		"Name":      deployment.ObjectMeta.Name,
		"Namespace": deployment.ObjectMeta.Namespace,
	}).Debugf("registering waiter for deployment")

	dwi.waitgroup.Add(1)
	go dwi.run(deployment)
	return nil
}

func (dwi *DeploymentWaitInterceptor) run(deployment *v1beta1.Deployment) {
	rs, err := kubeutil.GetReplicaSetForDeployment(dwi.client, deployment)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":      err,
			"StackTrace": fmt.Sprintf("%+v", err),
		}).Warn("Failed to get replica set for deployment")
		return
	}

	dwi.waitgroup.Add(1)
	go dwi.podNotifier(rs)

	selector := fields.OneTermEqualSelector("metadata.name", deployment.ObjectMeta.Name)
	for deployment := range kubeutil.WatchDeployments(dwi.ctx, dwi.client, selector) {
		if kubeutil.DeploymentRolloutComplete(deployment) {
			dwi.cancel()
		}
	}

	dwi.cancel()

	log.WithFields(log.Fields{
		"Name": deployment.ObjectMeta.Name,
	}).Debugf("deployment succeeded")

	dwi.waitgroup.Done()
}

func (dwi *DeploymentWaitInterceptor) podNotifier(rs *v1beta1.ReplicaSet) {
	for pod := range kubeutil.WatchPods(dwi.ctx, dwi.client, fields.Everything()) {
		if !kubeutil.IsOwner(rs.ObjectMeta, pod.ObjectMeta) {
			continue
		}

		log.WithFields(log.Fields{
			"Name":      pod.ObjectMeta.Name,
			"Namespace": pod.ObjectMeta.Namespace,
		}).Debugf("pod changed")

		err := kubeutil.PodWarnings(pod)
		if err != nil {
			log.WithFields(log.Fields{
				"Name":      pod.ObjectMeta.Name,
				"Namespace": pod.ObjectMeta.Namespace,
				"Status":    pod.Status,
			}).Warn(err)
		}
	}
	dwi.waitgroup.Done()
}
