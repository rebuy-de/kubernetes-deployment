package kubectl

import "k8s.io/apimachinery/pkg/runtime"

type Interface interface {
	Apply(obj runtime.Object) (runtime.Object, error)
}
