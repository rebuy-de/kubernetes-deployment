package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/git"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubernetes"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/util"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

type testKubectl struct {
	calls []string
}

func (k *testKubectl) Apply(manifestFile string) ([]byte, error) {
	log.Printf("$ kubectlMock apply %s", manifestFile)
	k.calls = append(k.calls, "apply "+path.Base(manifestFile))
	return []byte{}, nil
}

func (k *testKubectl) Get(manifestFile string) ([]byte, error) {
	log.Printf("$ kubectlMock get %s", manifestFile)
	return nil, fmt.Errorf("not implemented, yet")
}

func createTestDirs(t *testing.T, base string, dirs ...string) {
	for _, dir := range dirs {
		err := os.MkdirAll(path.Join(base, dir), 0755)
		if err != nil {
			t.Error("failed to create dir")
			t.Error(err)
			t.FailNow()
		}
	}
}

func createTestGitRepo(t *testing.T, repopath string, branch string, subpath string, files ...string) func() {
	util.AssertNoError(t, os.MkdirAll(repopath, 0755))

	git, err := git.New(repopath)
	util.AssertNoError(t, err)

	util.AssertNoError(t, git.Init())
	util.AssertNoError(t, git.Exec("config", "user.email", "me@example.com"))
	util.AssertNoError(t, git.Exec("config", "user.name", "git example user"))

	for _, file := range files {
		filepath := path.Join(repopath, subpath, file)
		util.AssertNoError(t, os.MkdirAll(path.Dir(filepath), 0755))
		util.AssertNoError(t,
			ioutil.WriteFile(filepath,
				[]byte("{{.clusterDomain}}"), 0644))

	}

	util.AssertNoError(t, git.Exec("add", "."))
	util.AssertNoError(t, git.Exec("commit", "-m", "initial master commit"))
	if branch != "master" {
		util.AssertNoError(t, git.Exec("checkout", "-b", branch))
	}

	return func() {
		os.RemoveAll(repopath)
	}
}

func prepareTestEnvironment(t *testing.T) (*App, *testKubectl, func()) {
	tempDir, cleanup := util.TestCreateTempDir(t)

	createTestGitRepo(t,
		path.Join(tempDir, "repos", "bish"),
		"master", "/deployment/k8s",
		"bish-a.yml", "bish-b.yml", "bish-c.yaml", "bish-d.txt", "foo/bish-e.yml")

	createTestGitRepo(t,
		path.Join(tempDir, "repos", "bash"),
		"special", "/deployment/k8s",
		"bash-a.yml", "bash-b.yml", "bash-c.yaml", "bash-d.txt", "foo/bash-e.yml")

	createTestGitRepo(t,
		path.Join(tempDir, "repos", "bosh"),
		"master", "/deployment/foo",
		"bosh-a.yml", "bosh-b.yml", "bosh-c.yaml", "bosh-d.txt", "foo/bosh-e.yml")

	config, err := settings.ReadProjectConfigFrom("test-fixtures/services_test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var finalConfig settings.ProjectConfig
	finalConfig.Settings = config.Settings
	finalConfig.Settings.Output = path.Join(tempDir, "output")
	finalConfig.Services = settings.Services{}

	for _, service := range config.Services {
		service.Repository = path.Join(tempDir, service.Repository)
		var serviceInstance = *service
		finalConfig.Services = append(finalConfig.Services, &serviceInstance)
	}

	finalConfig.WriteTo(path.Join(tempDir, "config.yml"))

	kubectlMock := new(testKubectl)

	return &App{
		KubectlBuilder: func(*string) (kubernetes.API, error) {
			return kubectlMock, nil
		},
		Config: finalConfig,

		IgnoreDeployFailures: false,

		Goals: []string{"all"},
	}, kubectlMock, cleanup
}

func TestSkipAll(t *testing.T) {
	var err error

	app, _, cleanup := prepareTestEnvironment(t)
	defer cleanup()

	app.IgnoreDeployFailures = false
	app.Goals = []string{"all"}
	if err != nil {
		t.Fatal(err)
	}

	err = app.PrepareConfig()
	if err != nil {
		t.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		t.Fatal(err)
	}

	if len(app.Config.Services) != 3 {
		t.Errorf("The generated config looks wrong. Expected 3 services, but got %d.", len(app.Config.Services))
	}
}

func testInArray(array []string, search string) bool {
	for _, item := range array {
		if item == search {
			return true
		}
	}
	return false
}

func TestWholeApplication(t *testing.T) {
	app, kubectlMock, cleanup := prepareTestEnvironment(t)
	defer cleanup()

	err := app.PrepareConfig()
	if err != nil {
		t.Fatal(err)
	}

	err = app.Run()
	if err != nil {
		t.Fatal(err)
	}

	manifests := map[string][]string{
		"bish": {
			"bish-a.yml",
			"bish-b.yml",
			"bish-c.yaml",
		},
		"bash": {
			"bash-a.yml",
			"bash-b.yml",
			"bash-c.yaml",
		},
		"bosh": {
			"bosh-a.yml",
			"bosh-b.yml",
			"bosh-c.yaml",
		},
	}

	totalFiles := 0
	templateContents := "{{.clusterDomain}}"
	renderedContents := "unit-test.example.org"

	for project, files := range manifests {
		for _, file := range files {
			totalFiles += 1

			template, err := ioutil.ReadFile(path.Join(app.Config.Settings.Output,
				"templates", project, file))
			if err != nil {
				t.Fatal(err)
			}
			if string(template) != templateContents {
				t.Errorf("Template %s/%s has wrong contents.", project, file)
				t.Errorf("  Expected %#v", templateContents)
				t.Errorf("  Obtained %#v", string(template))
			}

			rendered, err := ioutil.ReadFile(path.Join(app.Config.Settings.Output,
				"rendered", project, file))
			if err != nil {
				t.Fatal(err)
			}
			if string(rendered) != renderedContents {
				t.Errorf("Template %s/%s has wrong contents.", project, file)
				t.Errorf("  Expected %#v", renderedContents)
				t.Errorf("  Obtained %#v", string(rendered))
			}

			search := fmt.Sprintf("apply %s", file)
			if !testInArray(kubectlMock.calls, search) {
				t.Errorf("Missing kubectl call.")
				t.Errorf("  Expected %#v", search)
				t.Errorf("  Obtained %#v", kubectlMock.calls)
			}
		}
	}

	if totalFiles != len(kubectlMock.calls) {
		t.Errorf("kubectl has wrong call count.")
		t.Errorf("  Expected %d calls", totalFiles)
		t.Errorf("  Obtained %#v", kubectlMock.calls)
	}
}
