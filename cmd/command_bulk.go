package cmd

import (
	"log"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubernetes"
	"github.com/spf13/cobra"
)

var (
	defaultProjectConfigPath = "config/services.yaml"
)

func NewBulkCommand(app *App) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "bulk",
		Short: "deploys all services from config file",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		if len(app.Goals) == 0 {
			log.Fatalf("You have to specify at least one of these goals: %v", GoalOrder)
		}

		app.KubectlBuilder = func(kubeconfig *string) (kubernetes.API, error) {
			return kubernetes.New(*kubeconfig)
		}

		Must(app.PrepareConfig())
		Must(app.Run())
	}

	return cmd

}
