package api

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

func Apply(params *Parameters, project, branchName string) error {
	objects, err := Generate(params, project, branchName)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, obj := range objects {
		raw, err := json.MarshalIndent(obj, "", "    ")
		if err != nil {
			return errors.WithStack(err)
		}

		err = params.Kubectl().Apply(bytes.NewBuffer(raw))
		if err != nil {
			return errors.Wrap(err, "unable to apply manifest")
		}
	}

	return nil
}
