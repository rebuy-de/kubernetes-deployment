package cmd

import (
	"fmt"
	"math/rand"
)

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

func randomID() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, 7)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
