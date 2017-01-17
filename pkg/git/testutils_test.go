package git

import (
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/util"
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
