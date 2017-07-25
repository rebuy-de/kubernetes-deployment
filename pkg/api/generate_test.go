package api_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	fakeGH "github.com/rebuy-de/kubernetes-deployment/pkg/gh/fake"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
	"github.com/rebuy-de/kubernetes-deployment/pkg/testutil"

	log "github.com/sirupsen/logrus"
)

var (
	ExampleSettings *settings.Settings
	ExampleGitHub   = &fakeGH.GitHub{
		"rebuy-de": fakeGH.Repos{
			"foobar": fakeGH.Branches{
				"master": fakeGH.Branch{
					Meta: gh.Branch{
						Name: "master",
						SHA:  "1234567",
					},
				},
				"1234567": fakeGH.Branch{
					Files: fakeGH.ScanFiles("test-fixtures/repos/foobar/master"),
				},
			},
		},
	}
)

func TestMain(m *testing.M) {
	var err error
	ExampleSettings, err = settings.ReadFromFile("test-fixtures/deployments.yaml")
	if err != nil {
		panic(err)
	}

	flag.Parse()

	if testing.Verbose() {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}

	os.Exit(m.Run())
}

func generateApp() *api.App {
	return &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
			Statsd: statsdw.NullClient{},
		},
		Settings:     ExampleSettings,
		Parameters:   new(api.Parameters),
		Interceptors: interceptors.New(),
	}
}

func TestProjectNoExist(t *testing.T) {
	app := generateApp()

	_, err := app.Generate("project-no-exist", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "failed to fetch project: project 'project-no-exist' not found"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestMissingRepo(t *testing.T) {
	app := generateApp()

	_, err := app.Generate("repo-no-exist", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "failed to fetch project: unable to get branch information: fake repo 'rebuy-de/repo-no-exist' doesn't exist"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestMissingBranch(t *testing.T) {
	app := generateApp()

	_, err := app.Generate("foobar", "missing-branch")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "failed to fetch project: unable to get branch information: fake branch 'rebuy-de/foobar#missing-branch' doesn't exist"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestMissingFiles(t *testing.T) {
	app := generateApp()

	_, err := app.Generate("no-files", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "didn't find any template files"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestInvalidFile(t *testing.T) {
	app := generateApp()

	_, err := app.Generate("invalid-file", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := `failed to render project: unable to decode file 'invalid.yaml': couldn't get version/kind; json parse error: json: cannot unmarshal array into Go value of type struct { APIVersion string "json:\"apiVersion,omitempty\""; Kind string "json:\"kind,omitempty\"" }`
	if !strings.HasPrefix(err.Error(), expect) {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected prefix:  %s", expect)
		t.Errorf("  Obtained message: %v", err)
	}
}

func TestGenerateSuccess(t *testing.T) {
	app := generateApp()

	objects, err := app.Generate("foobar", "master")
	if err != nil {
		t.Fatal(err)
	}

	testutil.AssertGoldenJSON(t, "test-fixtures/generate-success-golden.json", objects)
}
