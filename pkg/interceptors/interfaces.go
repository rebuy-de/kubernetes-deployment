package interceptors

import "k8s.io/apimachinery/pkg/runtime"

type Interface interface {
	ManifestApplied
	AllManifestsApplied
	Closer
}

type ManifestApplied interface {
	ManifestApplied(runtime.Object) error
}

type AllManifestsApplied interface {
	AllManifestsApplied([]runtime.Object) error
}

type ManifestRendered interface {
	ManifestRendered(runtime.Object) (runtime.Object, error)
}

type Closer interface {
	Close() error
}
