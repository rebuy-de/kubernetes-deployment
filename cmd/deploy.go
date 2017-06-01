package cmd

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
)

func NewDeployCommand(params *api.Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy PROJECT [BRANCH]",
		Short: "Deploys a project to Kubernetes",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		project, branch, err := getProject(args)
		if err != nil {
			return err
		}

		api.Deploy(params, project, branch)
		return nil
	}

	return cmd
}
