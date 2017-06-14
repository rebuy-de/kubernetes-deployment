package cmd

import "testing"

func TestGetProject(t *testing.T) {
	var (
		project string
		branch  string
		err     error
	)

	project, branch, err = getProject([]string{})
	if err == nil {
		t.Errorf("Expected error on 0 args")
	}

	project, branch, err = getProject([]string{"blubber"})
	if err != nil {
		t.Errorf("Got error on single arg: %v", err)
	}
	if project != "blubber" {
		t.Errorf("Got wrong project on single arg: %s != blubber", project)
	}
	if branch != "master" {
		t.Errorf("Got wrong branch on single arg: %s != master", branch)
	}

	project, branch, err = getProject([]string{"bim", "baz"})
	if err != nil {
		t.Errorf("Got error on double arg: %v", err)
	}
	if project != "bim" {
		t.Errorf("Got wrong project on double arg: %s != bim", project)
	}
	if branch != "baz" {
		t.Errorf("Got wrong branch on single arg: %s != baz", branch)
	}

	project, branch, err = getProject([]string{"bish", "bash", "bosh"})
	if err == nil {
		t.Errorf("Expected error on 3 args")
	}
}
