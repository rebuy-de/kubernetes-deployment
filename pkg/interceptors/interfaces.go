package interceptors

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"k8s.io/apimachinery/pkg/runtime"
)

type Interface interface {
	PostFetcher
	PreApplier
	PreManifestApplier
	PostManifestApplier
	PostApplier
	PostManifestRenderer
	Closer
}

type PostFetcher interface {
	PostFetch(*gh.Branch) error
}

type PreApplier interface {
	PreApply([]runtime.Object) error
}

type PreManifestApplier interface {
	PreManifestApply(runtime.Object) (runtime.Object, error)
}

type PostManifestApplier interface {
	PostManifestApply(runtime.Object) error
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
