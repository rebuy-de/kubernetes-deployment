package statuschecker

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
)

func TestTypePreApplier(t *testing.T) {
	var inter interceptors.PreApplier
	inter = New(nil, "", "")
	_ = inter
}

func TestTypePostFetcher(t *testing.T) {
	var inter interceptors.PostFetcher
	inter = New(nil, "", "")
	_ = inter
}
