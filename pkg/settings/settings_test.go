package settings

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/testutil"
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
