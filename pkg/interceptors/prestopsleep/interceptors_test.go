package prestopsleep

import (
	"testing"

	v1apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/rebuy-go-sdk/testutil"
)

func TestType(t *testing.T) {
	var inter interceptors.PostManifestRenderer
	inter = New(5)
	_ = inter
}

func TestModify(t *testing.T) {
	deployment := &v1apps.Deployment{
		Spec: v1apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					InitContainers: []v1.Container{
						v1.Container{},
					},
					Containers: []v1.Container{
						v1.Container{},
						v1.Container{},
					},
				},
			},
		},
	}

	intercepted, err := New(5).PostManifestRender(deployment)
	if err != nil {
		t.Error(err)
		return
	}

	testutil.AssertGoldenJSON(t, "test-fixtures/deployment-golden.json", intercepted)
}
