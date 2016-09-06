package util

import (
	"io/ioutil"
	"os"
	"testing"
	"strings"
)

func TestCreateTempDir(t *testing.T) (string, func()) {
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

func AssertDirExists(t *testing.T, path string) {
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

func AssertFileExists(t *testing.T, path string) {
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

func AssertFileNotExists(t *testing.T, path string) {
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

func AssertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error("got unexpected error")
		t.Error(err)
		t.FailNow()
	}
}

func AssertHasError(t *testing.T, err error) {
	if err == nil {
		t.Error("Did not get any error, and was expecting one!")
		t.Error(err)
		t.FailNow()
	}
}

func AssertStringContains(t *testing.T, haystack string, needle string, msg string) {
	if !strings.Contains(haystack, needle) {
		t.Error("String dose not contain substring:", needle)
		if msg != "" {
			t.Error(msg)
		}
		t.FailNow()
	}
}

func AssertStringIsNotEmpty(t *testing.T, haystack string, msg string) {
	if haystack == "" {
		t.Error("String is empty!")
		if msg != "" {
			t.Error(msg)
		}
		t.FailNow()
	}
}