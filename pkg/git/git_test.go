package git

import (
	"io/ioutil"
	"os/exec"
	"path"
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/util"
)

func TestGitExec(t *testing.T) {
	git, cleanup := testCreateTempGitDir(t)
	defer cleanup()

	err := git.Exec("init")
	if err != nil {
		t.Error("failed to execute git")
		t.Error(err)
		t.FailNow()
	}

	util.AssertDirExists(t, path.Join(git.Directory, ".git"))
}

func TestGitCommitID(t *testing.T) {
	git, cleanup := testCreateTempGitDir(t)
	defer cleanup()

	must := func(msg string, err error) {
		if err != nil {
			t.Error(msg)
			t.Error(err)
			t.FailNow()
		}
	}

	must("git init failed", git.Init())
	must("git config failed", git.Config("user.email", "me@example.com"))
	must("git config failed", git.Config("user.name", "git example user"))
	must("create file failed", ioutil.WriteFile(path.Join(git.Directory, "test"), []byte("blubber"), 0644))
	must("git add failed", git.Exec("add", "."))
	must("git commit failed", git.Exec("commit", "-m", "test"))

	id, err := git.CommitID()
	must("getting hash", err)

	t.Logf("Commit ID: %s", id)
}

func TestGitExecWrongCommand(t *testing.T) {
	git, cleanup := testCreateTempGitDir(t)
	defer cleanup()

	err := git.Exec("not-existing-command")
	if err == nil {
		t.Error("expected error didn't occur")
		t.FailNow()
	}

	_, ok := err.(*exec.ExitError)
	if !ok {
		t.Error("got the wrong type of error")
		t.Error(err)
		t.FailNow()
	}
}

func TestGitInit(t *testing.T) {
	git, cleanup := testCreateTempGitDir(t)
	defer cleanup()

	err := git.Init()
	if err != nil {
		t.Error("failed to execute git")
		t.Error(err)
		t.FailNow()
	}

	util.AssertDirExists(t, path.Join(git.Directory, ".git"))
}
