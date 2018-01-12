package prestopsleep

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/testutil"
	"k8s.io/api/core/v1"
	v1beta1extensions "k8s.io/api/extensions/v1beta1"
)

func TestType(t *testing.T) {
	var inter interceptors.PostManifestRenderer
	inter = New(5)
	_ = inter
}

func TestModify(t *testing.T) {
	deployment := &v1beta1extensions.Deployment{
		Spec: v1beta1extensions.DeploymentSpec{
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
	New(5).AddToDeployment(deployment)

	testutil.AssertGoldenJSON(t, "test-fixtures/deployment-golden.json", deployment)
}
