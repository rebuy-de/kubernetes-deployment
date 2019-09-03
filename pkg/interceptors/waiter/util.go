package waiter

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func isDeployment(obj runtime.Object) bool {
	if obj == nil {
		return false
	}
	if obj.GetObjectKind() == nil {
		return false
	}

	return obj.GetObjectKind().GroupVersionKind().Kind == "Deployment"
}
