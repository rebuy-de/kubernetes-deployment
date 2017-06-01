package api

import (
	"fmt"

	"github.com/pkg/errors"
)

func Deploy(params *Parameters, project, branchName string) error {
	objects, err := Generate(params, project, branchName)
	if err != nil {
		return errors.WithStack(err)
	}

	fmt.Printf("%#v\n", objects)

	return nil
}
