package annotater

import (
	"testing"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/testutil"
	v1beta1extensions "k8s.io/api/extensions/v1beta1"
)

func TestTypePostManifestRenderer(t *testing.T) {
	var inter interceptors.PostManifestRenderer
	inter = New()
	_ = inter
}

func TestTypePostFetcher(t *testing.T) {
	var inter interceptors.PostFetcher
	inter = New()
	_ = inter
}

func TestModify(t *testing.T) {
	deployment := &v1beta1extensions.Deployment{}

	inter := New()

	err := inter.PostFetch(&gh.Branch{
		Author:  "bim baz",
		Date:    time.Unix(1234567890, 123456789).UTC(),
		Message: "fancy feature",
		SHA:     "1234567890abcdef",
		Location: gh.Location{
			Owner: "rebuy-de",
			Repo:  "example-silo",
			Path:  "deployment/k8s",
			Ref:   "master",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	obj, err := inter.PostManifestRender(deployment)
	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertGoldenJSON(t, "test-fixtures/deployment-golden.json", obj)
}
