package waiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	log "github.com/sirupsen/logrus"
	apps "k8s.io/api/apps/v1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

var (
	ErrImagePullMuteDuration = 30 * time.Second
)

type DeploymentWaitInterceptor struct {
	client kubernetes.Interface

	ctx              context.Context
	cancel           context.CancelFunc
	waitgroup        *sync.WaitGroup
	errImagePullMute time.Time
}

func NewDeploymentWaitInterceptor(client kubernetes.Interface) *DeploymentWaitInterceptor {
	ctx, cancel := context.WithCancel(context.Background())

	return &DeploymentWaitInterceptor{
		client:           client,
		ctx:              ctx,
		cancel:           cancel,
		waitgroup:        new(sync.WaitGroup),
		errImagePullMute: time.Now(),
	}
}

func (dwi *DeploymentWaitInterceptor) PostApply([]runtime.Object) error {
	dwi.waitgroup.Wait()
	dwi.cancel()
	return nil
}

func (dwi *DeploymentWaitInterceptor) Close() error {
	dwi.cancel()
	return nil
}

func (dwi *DeploymentWaitInterceptor) PostManifestApply(obj runtime.Object) error {
	if !isDeployment(obj) {
		return nil
	}

	metas := kubeutil.SubObjectAccessor(obj)
	if len(metas) == 0 {
		return nil
	}

	meta := metas[0]

	log.WithFields(log.Fields{
		"Namespace": meta.GetNamespace(),
		"Name":      meta.GetName(),
	}).Debugf("registering waiter for deployment")

	dwi.waitgroup.Add(1)
	go dwi.run(meta.GetNamespace(), meta.GetName())
	return nil
}

func (dwi *DeploymentWaitInterceptor) run(namespace, name string) {
	defer dwi.waitgroup.Done()

	ctx, done := context.WithCancel(dwi.ctx)
	defer done()

	// We need to sleep a short time to let the controller update the revision
	// number and then update the deployment to see the current revision.
	time.Sleep(1 * time.Second)
	deployment, err := dwi.client.
		AppsV1().
		Deployments(namespace).
		Get(name, v1meta.GetOptions{})

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

func (dwi *DeploymentWaitInterceptor) podNotifier(ctx context.Context, rs *apps.ReplicaSet) {
	defer dwi.waitgroup.Done()

	for pod := range kubeutil.WatchPods(ctx, dwi.client, fields.Everything()) {
		if !kubeutil.IsOwner(rs.ObjectMeta, pod.ObjectMeta) {
			continue
		}

		log.WithFields(log.Fields{
			"Manifest": pod,
		}).Debugf("pod changed")

		err := kubeutil.PodWarnings(pod)

		_, ok := err.(kubeutil.ErrImagePull)
		if ok {
			if time.Now().Before(dwi.errImagePullMute) {
				continue
			}
			dwi.errImagePullMute = time.Now().Add(ErrImagePullMuteDuration)
		}

		if err != nil {
			log.WithFields(log.Fields{
				"Name":      pod.ObjectMeta.Name,
				"Namespace": pod.ObjectMeta.Namespace,
				"PodData":   pod,
			}).Warn(err)
		}
	}
}
