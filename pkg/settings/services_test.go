package settings

import (
	"reflect"
	"testing"
)

func TestReadFile(t *testing.T) {
	settings, err := Read("./test-fixtures/services_test.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	expect := Settings{
		TemplateValues: TemplateValues{
			TemplateValue{
				Name:  "clusterDomain",
				Value: "unit-test.example.org",
			},
		},
		Services: Services{
			Service{
				Repository: "repos/bish",
			},
			Service{
				Repository: "repos/bash",
				Branch:     "special",
			},
			Service{
				Repository: "repos/bosh",
				Path:       "//deployment/foo",
			},
		},
	}

	if !reflect.DeepEqual(*settings, expect) {
		t.Errorf("Generated value doesn't match expected one.")
		t.Errorf("  Expected: %+v", expect)
		t.Errorf("  Obtained: %+v", settings)
	}

}
