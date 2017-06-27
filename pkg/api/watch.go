package api

import (
	"context"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
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

		rs, err := kubeutil.GetReplicaSetForDeployment(app.Clients.Kubernetes, &deployment)
		if err != nil {
			return errors.WithStack(err)
		}

		ctx, done := context.WithCancel(context.Background())

		go func() {
			for pod := range kubeutil.WatchPods(
				ctx, app.Clients.Kubernetes,
				fields.Everything()) {

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
		}()

		for deployment := range kubeutil.WatchDeployments(
			ctx, app.Clients.Kubernetes,
			fields.OneTermEqualSelector("metadata.name", deployment.ObjectMeta.Name)) {
			if kubeutil.DeploymentRolloutComplete(deployment) {
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
	}).Debugf("all deployments succeeded")

	return nil
}
