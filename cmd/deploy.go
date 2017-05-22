package cmd

import (
	"fmt"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
)

const (
	DefaultBranch = "master"
)

func NewDeployCommand(params *api.Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy PROJECT [BRANCH]",
		Short: "Deploys a project to Kubernetes",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("Wrong number of arguments.")
		}

		var (
			project = args[0]
			branch  = DefaultBranch
		)

		if len(args) > 1 {
			branch = args[1]
		}

		api.Deploy(params, project, branch)
		return nil
	}

	return cmd
}
