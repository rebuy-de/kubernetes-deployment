package fake

import (
	"reflect"
	"testing"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

var (
	ExampleBranch = gh.Branch{
		SHA:     "1234abc",
		Author:  "test",
		Date:    time.Now(),
		Message: "dis is a test",
	}

	Example gh.Interface = &GitHub{
		"rebuy-de": Repos{
			"info": Branches{
				"master": Branch{
					Meta: ExampleBranch,
					Files: Files{
						{Location: &gh.Location{Path: "deployments.yaml"}, Content: YAML([]string{"foo", "bar"})},
						{Location: &gh.Location{Path: "README.md"}, Content: "blubber"},
						{Location: &gh.Location{Path: "sub/foo.txt"}, Content: "bar"},
						{Location: &gh.Location{Path: "sub/bim.txt"}, Content: "baz"},
					},
				},
			},
		},
	}

	ExampleFile = &gh.Location{
		Owner: "rebuy-de",
		Repo:  "info",
		Ref:   "master",
		Path:  "deployments.yaml",
	}

	ExampleDir = &gh.Location{
		Owner: "rebuy-de",
		Repo:  "info",
		Ref:   "master",
		Path:  "/",
	}

	ExampleSubDir = &gh.Location{
		Owner: "rebuy-de",
		Repo:  "info",
		Ref:   "master",
		Path:  "sub",
	}
)

func TestGetBranch(t *testing.T) {
	branch, err := Example.GetBranch(ExampleFile)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ExampleBranch, *branch) {
		t.Errorf("Branch data doesn't match:")
		t.Errorf("  Expected: %#v", ExampleBranch)
		t.Errorf("  Obtained: %#v", *branch)
	}
}

func TestGetFile(t *testing.T) {
	file, err := Example.GetFile(ExampleFile)
	if err != nil {
		t.Fatal(err)
	}

	expected := "- foo\n- bar\n"

	if file.Content != expected {
		t.Errorf("File contents don't match:")
		t.Errorf("  Expected: %#v", expected)
		t.Errorf("  Obtained: %#v", file)
	}
}

func TestGetFiles(t *testing.T) {
	files, err := Example.GetFiles(ExampleDir)
	if err != nil {
		t.Fatal(err)
	}

	expected := []gh.File{
		{Location: &gh.Location{Path: "deployments.yaml"}, Content: "- foo\n- bar\n"},
		{Location: &gh.Location{Path: "README.md"}, Content: "blubber"},
	}

	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Files don't match:")
		t.Errorf("  Expected: %#v", expected)
		t.Errorf("  Obtained: %#v", files)
	}
}

func TestGetSubdirectoryFiles(t *testing.T) {
	files, err := Example.GetFiles(ExampleSubDir)
	if err != nil {
		t.Fatal(err)
	}

	expected := []gh.File{
		{Location: &gh.Location{Path: "sub/foo.txt"}, Content: "bar"},
		{Location: &gh.Location{Path: "sub/bim.txt"}, Content: "baz"},
	}

	if !reflect.DeepEqual(files, expected) {
		t.Errorf("Files don't match:")
		t.Errorf("  Expected: %#v", expected)
		t.Errorf("  Obtained: %#v", files)
	}
}
