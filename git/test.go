package git

import (
	"io/ioutil"
	"os"
	"testing"
)

func testCreateTempDir(t *testing.T) (string, func()) {
	tempDir, err := ioutil.TempDir("", "golang-git-test")
	if err != nil {
		t.Error("failed to create temporay directory")
		t.Error(err)
		t.FailNow()
	}
	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}

func testCreateTempGitDir(t *testing.T) (*Git, func()) {
	tempDir, tempDirDefer := testCreateTempDir(t)

	git, err := New(tempDir)
	if err != nil {
		t.Error("failed to initialize git")
		t.Error(err)
		t.FailNow()
	}

	return git, tempDirDefer
}

func assertDirExists(t *testing.T, path string) {
	info, err := os.Stat(path)

	if err != nil {
		t.Error("required directory doesn't exist")
		t.Error(err)
		t.FailNow()
	}

	if !info.IsDir() {
		t.Error("required path isn't a directory")
		t.FailNow()
	}
}

func assertFileExists(t *testing.T, path string) {
	info, err := os.Stat(path)

	if err != nil {
		t.Error("required file doesn't exist")
		t.Error(err)
		t.FailNow()
	}

	if info.IsDir() {
		t.Error("required path is a directory")
		t.FailNow()
	}
}

func assertFileNotExists(t *testing.T, path string) {
	_, err := os.Stat(path)

	if err == nil {
		t.Errorf("file '%s' exists", path)
		t.FailNow()
	}

	if !os.IsNotExist(err) {
		t.Error("got unexpected error")
		t.Error(err)
		t.FailNow()
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error("got unexpected error")
		t.Error(err)
		t.FailNow()
	}
}
