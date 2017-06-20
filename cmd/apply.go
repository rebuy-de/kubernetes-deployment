package cmd

import (
	"log"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
)

func NewApplyCommand(params *api.Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply PROJECT [BRANCH]",
		Short: "Deploys a project to Kubernetes",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		project, branch, err := getProject(args)
		if err != nil {
			return err
		}

		app, err := api.New(params)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		err = app.Apply(project, branch)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		return nil
	}

	return cmd
}
