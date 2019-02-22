package annotater

import (
	"fmt"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	batch_v1beta1 "k8s.io/api/batch/v1beta1"
	core_v1 "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
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
			obj: &apps_v1.Deployment{
				ObjectMeta: meta.ObjectMeta{Name: "deployment"},
			},
		},
		{
			name: "statefulset",
			obj: &apps_v1.StatefulSet{
				ObjectMeta: meta.ObjectMeta{Name: "statefulset"},
			},
		},
		{
			name: "daemonset",
			obj: &apps_v1.DaemonSet{
				ObjectMeta: meta.ObjectMeta{Name: "daemonset"},
			},
		},
		{
			name: "cronjob",
			obj: &batch_v1beta1.CronJob{
				ObjectMeta: meta.ObjectMeta{Name: "cronjob"},
			},
		},
		{
			name: "job",
			obj: &batch_v1.Job{
				ObjectMeta: meta.ObjectMeta{Name: "job"},
			},
		},
		{
			name: "service",
			obj: &core_v1.Service{
				ObjectMeta: meta.ObjectMeta{Name: "service"},
			},
		},
		{
			name: "pvc",
			obj: &core_v1.PersistentVolumeClaim{
				ObjectMeta: meta.ObjectMeta{
					Name: "pvc",
					Annotations: map[string]string{
						"volume.beta.kubernetes.io/storage-class": "aws-ebs-gp2",
					},
				},
			},
		},
	}

	mockClock := clock.NewMock()
	mockClock.Set(time.Unix(1234567890, 0).UTC())

	inter := New()
	inter.clock = mockClock
	inter.timezone = time.UTC

	err := inter.PostFetch(&gh.Branch{
		Name:    "master",
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
