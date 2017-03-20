package settings

import (
	"path"
	"reflect"
	"testing"

	"github.com/spf13/pflag"
)

func TestBuilder(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.PanicOnError)
	builder := NewBuilder(fs)
	fs.Set("config", path.Join("test-fixtures", "services_test.yaml"))

	config := builder()

	expect := ProjectConfig{
		Services: Services{
			&Service{
				Repository: "repos/bish",
			},
			&Service{
				Repository: "repos/bash",
				Branch:     "special",
			},
			&Service{
				Repository: "repos/bosh",
				Path:       "//deployment/foo",
			},
		},
		Settings: Settings{
			Kubeconfig:           "test-fixtures/kubeconfig.yml",
			Output:               "target/test-output",
			Sleep:                42000000000,
			SkipShuffle:          false,
			RetrySleep:           250000000,
			RetryCount:           3,
			IgnoreDeployFailures: false,
			TemplateValues: map[string]string{
				"cluster-domain": "unit-test.example.org",
			},
		},
	}

	if !reflect.DeepEqual(config, expect) {
		t.Errorf("Read config doesn't equal expectations:")
		t.Errorf("  Expected: %#v", expect)
		t.Errorf("  Obtained: %#v", config)
	}

}
