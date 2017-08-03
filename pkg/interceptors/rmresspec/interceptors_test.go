package rmresspec

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
)

func TestType(t *testing.T) {
	var inter interceptors.PostManifestRenderer
	inter = New()
	_ = inter
}
