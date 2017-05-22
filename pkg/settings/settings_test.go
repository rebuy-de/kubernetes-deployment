package settings

import (
	"reflect"
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

func TestReadFile(t *testing.T) {
	settings, err := Read("./test-fixtures/services_test.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	expect := Settings{
		Default: Defaults{
			Location: gh.Location{
				Owner: "rebuy-de",
				Path:  "deployment/k8s/",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "unit-test.example.org",
				"secret":        "foo",
			},
		},
		Services: Services{
			Service{
				Name: "",
				Location: gh.Location{
					Repo: "bish",
				},
			},
			Service{
				Name: "",
				Location: gh.Location{
					Repo: "bash",
				},
				TemplateValues: TemplateValues{
					"clusterDomain": "test.example.com",
				},
			},
			Service{
				Name: "",
				Location: gh.Location{
					Repo: "bosh", Path: "//deployment/foo",
				},
			},
			Service{
				Name: "foo",
				Location: gh.Location{
					Repo: "bar",
				},
			},
			Service{
				Name: "",
				Location: gh.Location{
					Owner: "kubernetes", Repo: "blub",
				},
			},
			Service{
				Location: gh.Location{
					Repo: "meh",
					Path: "deployment/k8s",
				},
			},
		},
	}

	if !reflect.DeepEqual(*settings, expect) {
		t.Errorf("Generated value doesn't match expected one.")
		t.Errorf("  Expected: %#v", expect)
		t.Errorf("  Obtained: %#v", *settings)
	}

}

func TestClean(t *testing.T) {
	settings, err := Read("./test-fixtures/services_test.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}

	settings.Clean()

	expect := Services{
		Service{
			Name: "bish",
			Location: gh.Location{
				Owner: "rebuy-de",
				Repo:  "bish",
				Path:  "deployment/k8s/",
				Ref:   "master",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "unit-test.example.org",
				"secret":        "foo",
			},
		},
		Service{
			Name: "bash",
			Location: gh.Location{
				Owner: "rebuy-de",
				Repo:  "bash",
				Path:  "deployment/k8s/",
				Ref:   "master",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "test.example.com",
				"secret":        "foo",
			},
		},
		Service{
			Name: "bosh-deployment-foo",
			Location: gh.Location{
				Owner: "rebuy-de",
				Repo:  "bosh",
				Path:  "deployment/foo/",
				Ref:   "master",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "unit-test.example.org",
				"secret":        "foo",
			},
		},
		Service{
			Name: "foo",
			Location: gh.Location{
				Owner: "rebuy-de",
				Repo:  "bar",
				Path:  "deployment/k8s/",
				Ref:   "master",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "unit-test.example.org",
				"secret":        "foo",
			},
		},
		Service{
			Name: "kubernetes-blub",
			Location: gh.Location{
				Owner: "kubernetes",
				Repo:  "blub",
				Path:  "deployment/k8s/",
				Ref:   "master",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "unit-test.example.org",
				"secret":        "foo",
			},
		},
		Service{
			Name: "meh",
			Location: gh.Location{
				Owner: "rebuy-de",
				Repo:  "meh",
				Path:  "deployment/k8s/",
				Ref:   "master",
			},
			TemplateValues: TemplateValues{
				"clusterDomain": "unit-test.example.org",
				"secret":        "foo",
			},
		},
	}

	if !reflect.DeepEqual(settings.Services, expect) {
		t.Errorf("Generated value doesn't match expected one.")
		t.Errorf("  Expected: %#v", expect)
		t.Errorf("  Obtained: %#v", settings.Services)
	}
}
