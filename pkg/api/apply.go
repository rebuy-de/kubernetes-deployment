package api

import "github.com/pkg/errors"

func (app *App) Apply(project, branchName string) error {
	objects, err := app.Generate(project, branchName)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, obj := range objects {
		err = app.Clients.Kubernetes.Apply(obj)
		if err != nil {
			return errors.Wrap(err, "unable to apply manifest")
		}
	}

	return nil
}
