package cmd

import "fmt"

const (
	DefaultBranch = "master"
)

func getProject(args []string) (string, string, error) {
	if len(args) < 1 || len(args) > 2 {
		return "", "", fmt.Errorf("Wrong number of arguments.")
	}

	var (
		project = args[0]
		branch  = DefaultBranch
	)

	if len(args) > 1 {
		branch = args[1]
	}

	return project, branch, nil

}
