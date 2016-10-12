package git

import (
	"github.com/rebuy-de/kubernetes-deployment/util"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func testCreateDirs(t *testing.T, base string, dirs ...string) {
	for _, dir := range dirs {
		err := os.MkdirAll(path.Join(base, dir), 0755)
		if err != nil {
			t.Error("failed to create dir")
			t.Error(err)
			t.FailNow()
		}
	}
}

func testTouchFiles(t *testing.T, base string, files ...string) {
	for _, file := range files {
		err := ioutil.WriteFile(path.Join(base, file),
			[]byte(time.Now().String()), 0644)
		if err != nil {
			t.Error("failed to touch file")
			t.Error(err)
			t.FailNow()
		}
	}
}

func TestSparseCheckout(t *testing.T) {
	git, cleanupParent := testCreateTempGitDir(t)
	defer cleanupParent()

	util.AssertNoError(t, git.Init())
	util.AssertNoError(t, git.Exec("config", "user.email", "me@example.com"))
	util.AssertNoError(t, git.Exec("config", "user.name", "git example user"))

	testCreateDirs(t, git.Directory,
		"foo/bar",
		"bim/baz",
		"bish/bash/bosh",
	)

	testTouchFiles(t, git.Directory,
		"foo/bar/a.yml",
		"bim/baz/b.txt",
		"bim/baz/c.yml",
	)

	util.AssertNoError(t, git.Exec("add", "."))
	util.AssertNoError(t, git.Exec("commit", "-m", "initial commit"))
	util.AssertNoError(t, git.Exec("checkout", "-b", "second-branch"))

	testTouchFiles(t, git.Directory,
		"foo/bar/d.txt",
		"bim/baz/e.yml",
		"bish/bash/bosh/f.yml",
	)

	util.AssertNoError(t, git.Exec("add", "."))
	util.AssertNoError(t, git.Exec("commit", "-m", "second commit"))

	targetMaster, cleanupTargetMaster := util.TestCreateTempDir(t)
	defer cleanupTargetMaster()

	SparseCheckout(targetMaster, git.Directory, "master", "foo/bar")
	util.AssertFileExists(t, path.Join(targetMaster, "foo/bar/a.yml"))
	util.AssertFileNotExists(t, path.Join(targetMaster, "bim/baz/b.txt"))
	util.AssertFileNotExists(t, path.Join(targetMaster, "bim/baz/c.yml"))
	util.AssertFileNotExists(t, path.Join(targetMaster, "foo/bar/d.txt"))
	util.AssertFileNotExists(t, path.Join(targetMaster, "bim/baz/e.yml"))
	util.AssertFileNotExists(t, path.Join(targetMaster, "bish/bash/bosh/f.yml"))

	targetBranch, cleanupTargetBranch := util.TestCreateTempDir(t)
	defer cleanupTargetBranch()

	SparseCheckout(targetBranch, git.Directory, "second-branch", "foo/bar")
	util.AssertFileExists(t, path.Join(targetBranch, "foo/bar/a.yml"))
	util.AssertFileNotExists(t, path.Join(targetBranch, "bim/baz/b.txt"))
	util.AssertFileNotExists(t, path.Join(targetBranch, "bim/baz/c.yml"))
	util.AssertFileExists(t, path.Join(targetBranch, "foo/bar/d.txt"))
	util.AssertFileNotExists(t, path.Join(targetBranch, "bim/baz/e.yml"))
	util.AssertFileNotExists(t, path.Join(targetBranch, "bish/bash/bosh/f.yml"))
}
