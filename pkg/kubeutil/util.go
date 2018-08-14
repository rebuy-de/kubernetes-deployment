package kubeutil

import (
	"github.com/pkg/errors"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
)

func IsOwner(parent, child v1meta.ObjectMeta) bool {
	for _, or := range child.OwnerReferences {
		if or.UID == parent.UID {
			return true
		}
	}
	return false
}

func Decode(raw []byte) (runtime.Object, error) {
	decode := scheme.Codecs.UniversalDeserializer().Decode
	obj, _, err := decode(raw, nil, nil)

	// Fallback to UnknownObject if the API/Kind is not registered, so the
	// interceptors still work. In case the Kind actually does not exist,
	// kubectl will fail later anyway.
	if runtime.IsNotRegisteredError(err) {
		unknown := new(UnknownObject)
		obj = unknown

		// Since JSON is a subset of YAML, this works for both.
		err = unknown.FromYAML(raw)
	}

	return obj, errors.Wrapf(err, "unable to decode manifest")
}
