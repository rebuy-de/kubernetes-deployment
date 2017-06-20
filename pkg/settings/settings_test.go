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

	testutil.AssertGoldenFile(t, "test-fixtures/services-plain-golden.json", settings)
}

func TestClean(t *testing.T) {
	settings, err := Read("./test-fixtures/services.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean()

	testutil.AssertGoldenFile(t, "test-fixtures/services-clean-golden.json", settings)
}
