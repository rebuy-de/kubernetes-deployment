package waiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	log "github.com/sirupsen/logrus"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	dwi.cancel()
	return nil
}

func (dwi *DeploymentWaitInterceptor) Close() error {
	dwi.cancel()
	return nil
}

func (dwi *DeploymentWaitInterceptor) ManifestApplied(obj runtime.Object) error {
	deployment, ok := obj.(*v1beta1.Deployment)
	if !ok {
		return nil
	}

	log.WithFields(log.Fields{
		"Manifest": deployment,
	}).Debugf("registering waiter for deployment")

	dwi.waitgroup.Add(1)
	go dwi.run(deployment)
	return nil
}

func (dwi *DeploymentWaitInterceptor) run(deployment *v1beta1.Deployment) {
	defer dwi.waitgroup.Done()

	ctx, done := context.WithCancel(dwi.ctx)
	defer done()

	// We need to sleep a short time to let the controller update the revision
	// number and then update the deployment to see the current revision.
	time.Sleep(1 * time.Second)
	deployment, err := dwi.client.
		Extensions().
		Deployments(deployment.ObjectMeta.Namespace).
		Get(deployment.ObjectMeta.Name, v1meta.GetOptions{})

	rs, err := kubeutil.GetReplicaSetForDeployment(dwi.client, deployment)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":      err,
			"StackTrace": fmt.Sprintf("%+v", err),
		}).Warn("Failed to get replica set for deployment")
		return
	}

	log.WithFields(log.Fields{
		"Name":      rs.ObjectMeta.Name,
		"Namespace": rs.ObjectMeta.Namespace,
	}).Debugf("found replica set for deployment")

	dwi.waitgroup.Add(1)
	go dwi.podNotifier(ctx, rs)

	selector := fields.OneTermEqualSelector("metadata.name", deployment.ObjectMeta.Name)
	for deployment := range kubeutil.WatchDeployments(ctx, dwi.client, selector) {
		if kubeutil.DeploymentRolloutComplete(deployment) {
			break
		}
	}

	log.WithFields(log.Fields{
		"Name":      deployment.ObjectMeta.Name,
		"Namespace": deployment.ObjectMeta.Namespace,
	}).Debugf("deployment succeeded")
}

func (dwi *DeploymentWaitInterceptor) podNotifier(ctx context.Context, rs *v1beta1.ReplicaSet) {
	defer dwi.waitgroup.Done()

	for pod := range kubeutil.WatchPods(ctx, dwi.client, fields.Everything()) {
		if !kubeutil.IsOwner(rs.ObjectMeta, pod.ObjectMeta) {
			continue
		}

		log.WithFields(log.Fields{
			"Manifest": pod,
		}).Debugf("pod changed")

		err := kubeutil.PodWarnings(pod)
		if err != nil {
			log.WithFields(log.Fields{
				"Name":      pod.ObjectMeta.Name,
				"Namespace": pod.ObjectMeta.Namespace,
			}).Warn(err)
		}
	}
}
