package rmoldjob

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
)

func TestTypePreManifestApplier(t *testing.T) {
	var inter interceptors.PreManifestApplier
	inter = &Interceptor{}
	_ = inter
}
