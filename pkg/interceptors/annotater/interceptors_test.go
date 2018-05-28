package annotater

import (
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	apps "k8s.io/api/apps/v1beta1"
	core "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/rebuy-go-sdk/testutil"
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
	cases := []struct {
		name string
		obj  runtime.Object
	}{
		{
			name: "deployment",
			obj:  &extensions.Deployment{},
		},
		{
			name: "statefulset",
			obj:  &apps.StatefulSet{},
		},
		{
			name: "service",
			obj:  &core.Service{},
		},
	}

	mockClock := clock.NewMock()
	mockClock.Set(time.Unix(1234567890, 0).UTC())

	inter := New()
	inter.clock = mockClock

	err := inter.PostFetch(&gh.Branch{
		Author:  "bim baz",
		Date:    mockClock.Now(),
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			obj, err := inter.PostManifestRender(tc.obj)
			if err != nil {
				t.Fatal(err)
			}

			golden := fmt.Sprintf("test-fixtures/%s-golden.json", tc.name)
			testutil.AssertGoldenJSON(t, golden, obj)
		})
	}
}
