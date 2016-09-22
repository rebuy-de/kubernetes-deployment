package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"testing"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/git"
	"github.com/rebuy-de/kubernetes-deployment/kubernetes"
	"github.com/rebuy-de/kubernetes-deployment/settings"
	"github.com/rebuy-de/kubernetes-deployment/util"
)

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
				[]byte(time.Now().String()), 0644))

	}

	util.AssertNoError(t, git.Exec("add", "."))
	if branch != "master" {
		util.AssertNoError(t, git.Exec("checkout", "-b", branch))
	}
	util.AssertNoError(t, git.Exec("commit", "-m", "initial commit"))

	return func() {
		os.RemoveAll(repopath)
	}
}

func prepareTestEnvironment(t *testing.T) (*App, *testKubectl, func()) {
	tempDir, cleanup := util.TestCreateTempDir(t)

	cleanup = createTestGitRepo(t,
		path.Join(tempDir, "repos", "bish"),
		"master", "/deployment/k8s",
		"bish-a.yml", "bish-b.yml", "bish-c.yaml", "bish-d.txt", "foo/bish-e.yml")

	cleanup = createTestGitRepo(t,
		path.Join(tempDir, "repos", "bash"),
		"special", "/deployment/k8s",
		"bash-a.yml", "bash-b.yml", "bash-c.yaml", "bash-d.txt", "foo/bash-e.yml")

	cleanup = createTestGitRepo(t,
		path.Join(tempDir, "repos", "bosh"),
		"master", "/deployment/foo",
		"bosh-a.yml", "bosh-b.yml", "bosh-c.yaml", "bosh-d.txt", "foo/bosh-e.yml")

	config, err := settings.ReadProjectConfigFrom("config/services_test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	var finalConfig settings.ProjectConfig
	finalConfig.Settings = config.Settings

	var finalServicesInstance settings.Services
	for _, service := range *config.Services {
		service.Repository = path.Join(tempDir, service.Repository)
		var serviceInstance = *service
		finalServicesInstance = append(finalServicesInstance, &serviceInstance)
	}

	finalConfig.Services = &finalServicesInstance

	finalConfig.WriteTo(path.Join(tempDir, "config.yml"))

	kubectlMock := new(testKubectl)

	return &App{
		KubectlBuilder: func(*string) (kubernetes.API, error) {
			return kubectlMock, nil
		},
		ProjectConfigPath: path.Join(tempDir, "config.yml"),
		OutputPath:        path.Join(tempDir, "output"),

		SleepInterval:        250 * time.Millisecond,
		IgnoreDeployFailures: false,

		RetrySleep: 250 * time.Millisecond,
		RetryCount: 3,

		SkipShuffle: false,
		SkipFetch:   false,
		SkipDeploy:  false,
	}, kubectlMock, cleanup
}

func TestSkipAll(t *testing.T) {
	app, _, cleanup := prepareTestEnvironment(t)
	defer cleanup()

	app.SkipShuffle = true
	app.SkipFetch = true
	app.SkipDeploy = true
	app.IgnoreDeployFailures = false

	err := app.Run()
	if err != nil {
		t.Fatal(err)
	}

	config, err := settings.ReadProjectConfigFrom(path.Join(app.OutputPath, "config.yml"))
	fmt.Println(config)
	util.AssertNoError(t, err)

	if len(*config.Services) != 3 {
		t.Errorf("The generated config looks wrong. Expected 3 services, but got %d.", len(*config.Services))
	}
}

func TestWholeApplication(t *testing.T) {
	app, kubectlMock, cleanup := prepareTestEnvironment(t)
	defer cleanup()

	err := app.Run()
	if err != nil {
		t.Fatal(err)
	}

	_, err = settings.ReadProjectConfigFrom(path.Join(app.OutputPath, "config.yml"))
	util.AssertNoError(t, err)

	calls := []string{
		"apply bish-a.yml",
		"apply bish-b.yml",
		"apply bish-c.yaml",
		"apply bash-a.yml",
		"apply bash-b.yml",
		"apply bash-c.yaml",
		"apply bosh-a.yml",
		"apply bosh-b.yml",
		"apply bosh-c.yaml",
	}

	sort.Strings(calls)
	sort.Strings(kubectlMock.calls)

	if !reflect.DeepEqual(calls, kubectlMock.calls) {
		t.Errorf("kubectl was called to often.")
		t.Errorf("  Expected %#v", calls)
		t.Errorf("  Obtained %#v", kubectlMock.calls)
	}

}
