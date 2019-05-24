package kubeutil

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/api/extensions/v1beta1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func DeploymentRolloutComplete(deployment *v1beta1.Deployment) bool {
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

func GetReplicaSetForDeployment(client kubernetes.Interface, deployment *v1beta1.Deployment) (*v1beta1.ReplicaSet, error) {
	replicaSets, err := client.
		ExtensionsV1beta1().
		ReplicaSets(deployment.ObjectMeta.Namespace).
		List(v1meta.ListOptions{})
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

		if IsOwner(deployment.ObjectMeta, rs.ObjectMeta) {
			return &rs, nil
		}
	}

	return nil, errors.Errorf("could not found replicaset for deployment")
}
