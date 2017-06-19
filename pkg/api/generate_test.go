package api_test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"k8s.io/client-go/pkg/api/v1"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	fakeGH "github.com/rebuy-de/kubernetes-deployment/pkg/gh/fake"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"

	log "github.com/Sirupsen/logrus"
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

func TestProjectNoExist(t *testing.T) {
	app := &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
		},
		Settings: ExampleSettings,
	}

	_, err := app.Generate("project-no-exist", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "project 'project-no-exist' not found"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestMissingRepo(t *testing.T) {
	app := &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
		},
		Settings: ExampleSettings,
	}

	_, err := app.Generate("repo-no-exist", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "Unable to get branch information: fake repo 'rebuy-de/repo-no-exist' doesn't exist"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestMissingBranch(t *testing.T) {
	app := &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
		},
		Settings: ExampleSettings,
	}

	_, err := app.Generate("foobar", "missing-branch")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "Unable to get branch information: fake branch 'rebuy-de/foobar#missing-branch' doesn't exist"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestMissingFiles(t *testing.T) {
	app := &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
		},
		Settings: ExampleSettings,
	}

	_, err := app.Generate("no-files", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "directory doesn't contain any template files"
	if err.Error() != expect {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected: %s", expect)
		t.Errorf("  Obtained: %v", err)
	}
}

func TestInvalidFile(t *testing.T) {
	app := &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
		},
		Settings: ExampleSettings,
	}

	_, err := app.Generate("invalid-file", "master")
	if err == nil {
		t.Fatal("expected an error")
	}

	expect := "unable to decode file 'invalid.yaml': "
	if !strings.HasPrefix(err.Error(), expect) {
		t.Errorf("Got wrong error:")
		t.Errorf("  Expected prefix:  %s", expect)
		t.Errorf("  Obtained message: %v", err)
	}
}

func TestGenerateSuccess(t *testing.T) {
	app := &api.App{
		Clients: &api.Clients{
			GitHub: ExampleGitHub,
		},
		Settings: ExampleSettings,
	}

	objects, err := app.Generate("foobar", "master")
	if err != nil {
		t.Fatal(err)
	}

	if len(objects) != 1 {
		t.Fatalf("Expected 1 object. Got %d.", len(objects))
	}

	pod := objects[0].(*v1.Pod)

	AssertGoldenFile(t, "test-fixtures/generate-success-golden.json", pod)
}
