package git

import (
	"github.com/rebuy-de/kubernetes-deployment/util"
	"testing"
)

func testCreateTempGitDir(t *testing.T) (*Git, func()) {
	tempDir, tempDirDefer := util.TestCreateTempDir(t)

	git, err := New(tempDir)
	if err != nil {
		t.Error("failed to initialize git")
		t.Error(err)
		t.FailNow()
	}

	return git, tempDirDefer
}
