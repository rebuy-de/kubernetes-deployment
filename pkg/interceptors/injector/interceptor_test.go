package injector

import (
	v1apps "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os/exec"
	"testing"

	"github.com/rebuy-de/rebuy-go-sdk/testutil"
)

func TestInterceptor_PostManifestRender(t *testing.T) {
	_, err := exec.LookPath("linkerd")

	if err != nil {
		t.Skipf("linkerd injector not tested, linkerd binary not present.")
	}

	deployment := &v1apps.Deployment{
		TypeMeta: meta.TypeMeta{
			APIVersion: "extensions/v1beta1",
			Kind:       "Deployment",
		},
		ObjectMeta: meta.ObjectMeta{
			Name: "linkerd-test",
			Annotations: map[string]string{
				"rebuy.com/kubernetes-deployment.inject-linkerd": "true",
			},
		},
		Spec: v1apps.DeploymentSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						v1.Container{
							Name: "container1",
						},
						v1.Container{
							Name: "container2",
						},
					},
				},
			},
		},
		Status: v1apps.DeploymentStatus{},
	}

	opts := Options{
		InjectArguments: []string{
			"--proxy-memory-request", "20Mi",
			"--proxy-cpu-request", "35m",
			"--ignore-cluster=true",
			"--linkerd-version=2.3.0",
			"--disable-identity",
		},
		ConnectTimeout: "10s",
	}

	intercepted, err := New(opts).PostManifestRender(deployment)
	if err != nil {
		t.Error(err)
		return
	}

	testutil.AssertGoldenJSON(t, "test-fixtures/deployment-golden.json", intercepted)
}
