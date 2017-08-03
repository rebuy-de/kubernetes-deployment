package interceptors

import "k8s.io/apimachinery/pkg/runtime"

type Interface interface {
	PreManifestApplier
	PostApplier
	PostManifestRenderer
	Closer
}

type PreManifestApplier interface {
	PreManifestApply(runtime.Object) error
}

type PostApplier interface {
	PostApply([]runtime.Object) error
}

type PostManifestRenderer interface {
	PostManifestRender(runtime.Object) (runtime.Object, error)
}

type Closer interface {
	Close() error
}
