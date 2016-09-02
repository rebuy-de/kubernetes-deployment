package git

import (
	"os/exec"
	"path"
	"testing"
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

	assertDirExists(t, path.Join(git.Directory, ".git"))
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

	assertDirExists(t, path.Join(git.Directory, ".git"))
}
