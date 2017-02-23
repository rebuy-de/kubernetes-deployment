package cmd

import (
	"log"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubernetes"
	"github.com/spf13/cobra"
)

func NewCmdDeploy(app *App) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploys a single project",
	}

	var branch string

	cmd.PersistentFlags().StringVarP(
		&branch,
		"branch", "b", "master",
		"branch to deploy")
	cmd.PersistentFlags().StringVarP(
		&app.Target,
		"name", "n", "",
		"name of the project to deploy")

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if len(app.Goals) == 0 {
			log.Fatalf("You have to specify at least one of these goals: %v", GoalOrder)
		}

		if len(app.Target) <= 0 {
			log.Fatalf("You have to specify the name of the project.")
		}

		app.KubectlBuilder = func(kubeconfig *string) (kubernetes.API, error) {
			return kubernetes.New(*kubeconfig)
		}

		Must(app.PrepareConfig())

		if len(branch) > 0 {
			app.Config.Services[0].Branch = branch
		}

		Must(app.Run())
	}

	return cmd

}
