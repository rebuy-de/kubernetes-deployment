package git

import (
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

	assertNoError(t, git.Init())
	assertNoError(t, git.Exec("config", "user.email", "me@example.com"))
	assertNoError(t, git.Exec("config", "user.name", "git example user"))

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

	assertNoError(t, git.Exec("add", "."))
	assertNoError(t, git.Exec("commit", "-m", "initial commit"))
	assertNoError(t, git.Exec("checkout", "-b", "second-branch"))

	testTouchFiles(t, git.Directory,
		"foo/bar/d.txt",
		"bim/baz/e.yml",
		"bish/bash/bosh/f.yml",
	)

	assertNoError(t, git.Exec("add", "."))
	assertNoError(t, git.Exec("commit", "-m", "second commit"))

	targetMaster, cleanupTargetMaster := testCreateTempDir(t)
	defer cleanupTargetMaster()

	SparseCheckout(targetMaster, git.Directory, "master", "foo/bar")
	assertFileExists(t, path.Join(targetMaster, "foo/bar/a.yml"))
	assertFileNotExists(t, path.Join(targetMaster, "bim/baz/b.txt"))
	assertFileNotExists(t, path.Join(targetMaster, "bim/baz/c.yml"))
	assertFileNotExists(t, path.Join(targetMaster, "foo/bar/d.txt"))
	assertFileNotExists(t, path.Join(targetMaster, "bim/baz/e.yml"))
	assertFileNotExists(t, path.Join(targetMaster, "bish/bash/bosh/f.yml"))

	targetBranch, cleanupTargetBranch := testCreateTempDir(t)
	defer cleanupTargetBranch()

	SparseCheckout(targetBranch, git.Directory, "second-branch", "foo/bar")
	assertFileExists(t, path.Join(targetBranch, "foo/bar/a.yml"))
	assertFileNotExists(t, path.Join(targetBranch, "bim/baz/b.txt"))
	assertFileNotExists(t, path.Join(targetBranch, "bim/baz/c.yml"))
	assertFileExists(t, path.Join(targetBranch, "foo/bar/d.txt"))
	assertFileNotExists(t, path.Join(targetBranch, "bim/baz/e.yml"))
	assertFileNotExists(t, path.Join(targetBranch, "bish/bash/bosh/f.yml"))
}
