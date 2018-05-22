package settings

import (
	"fmt"
	"testing"

	"github.com/rebuy-de/rebuy-go-sdk/testutil"
)

func TestReadFile(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertGoldenYAML(t, "test-fixtures/services-plain-golden.yaml", settings)
}

func TestClean(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean("")

	testutil.AssertGoldenYAML(t, "test-fixtures/services-clean-golden.yaml", settings)
}

func TestCleanWithContext(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean("def")

	testutil.AssertGoldenYAML(t, "test-fixtures/services-context-golden.yaml", settings)
}

func TestServiceGuessing(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil)
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
