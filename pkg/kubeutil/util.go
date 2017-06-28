package kubeutil

import (
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func IsOwner(parent, child v1meta.ObjectMeta) bool {
	for _, or := range child.OwnerReferences {
		if or.UID == parent.UID {
			return true
		}
	}
	return false
}
