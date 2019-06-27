package imagechecker

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
)

func TestTypePreApplier(t *testing.T) {
	var inter interceptors.PreApplier
	inter = New(nil, Options{})
	_ = inter
}
