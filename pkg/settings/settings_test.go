package settings

import (
	"fmt"
	"io/ioutil"
	"testing"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/rebuy-de/rebuy-go-sdk/testutil"
)

func TestReadFile(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertGoldenYAML(t, "test-fixtures/services-plain-golden.yaml", settings)
}

func TestReadConfigMap(t *testing.T) {
	original, err := ioutil.ReadFile("./test-fixtures/services.yaml")
	if err != nil {
		t.Fatal(err)
	}

	cm := &core.ConfigMap{
		ObjectMeta: meta.ObjectMeta{
			Name:      "kubernetes-deployment",
			Namespace: "default",
		},
		Data: map[string]string{
			"settings.yaml": string(original),
		},
	}

	kube := fake.NewSimpleClientset(cm)
	settings, err := Read("/api/v1/namespaces/default/configmaps/kubernetes-deployment", nil, kube)
	if err != nil {
		t.Fatal(err)
	}

	// This uses the same golden file as TestReadFile, but this is fine since
	// they actually should look the same.
	testutil.AssertGoldenYAML(t, "test-fixtures/services-plain-golden.yaml", settings)
}

func TestClean(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean("")

	testutil.AssertGoldenYAML(t, "test-fixtures/services-clean-golden.yaml", settings)
}

func TestCleanWithContext(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean("def")

	testutil.AssertGoldenYAML(t, "test-fixtures/services-context-golden.yaml", settings)
}

func TestServiceGuessing(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean("")

	cases := []struct {
		input  string
		result string
	}{
		{input: "my-service", result: "github.com/rebuy-de/my-service/deployment/k8s/@master"},
		{input: "my-service/sub", result: "github.com/rebuy-de/my-service/deployment/k8s/sub/@master"},
		{input: "bish", result: "github.com/rebuy-de/bish/deployment/k8s/@master"},
		{input: "guess", result: "github.com/rebuy-de/k8s-guess/other/@master"},
		{input: "guess/blub", result: "github.com/rebuy-de/k8s-guess/other/blub/@master"},
		{input: "cloud/prom", result: "github.com/rebuy-de/cloud/prom/@master"},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			svc := settings.Service(tc.input)

			want := tc.result
			have := fmt.Sprint(svc.Location)

			if want != have {
				t.Errorf("Wrong result.\n\tWant: %s.\n\tHave: %s.", want, have)
			}
		})
	}

}
