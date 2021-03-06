package waiter

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
)

func TestTypePostApplier(t *testing.T) {
	var inter interceptors.PostApplier
	inter = &DeploymentWaitInterceptor{}
	_ = inter
}

func TestTypePostManifestApplier(t *testing.T) {
	var inter interceptors.PostManifestApplier
	inter = &DeploymentWaitInterceptor{}
	_ = inter
}

func TestTypeCloser(t *testing.T) {
	var inter interceptors.Closer
	inter = &DeploymentWaitInterceptor{}
	_ = inter
}
