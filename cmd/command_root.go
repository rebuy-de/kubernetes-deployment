package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {

	var settingsBuilder settings.SettingsBuilder
	app := new(App)

	cmd := &cobra.Command{
		Use:   "kubernetes-deployment",
		Short: "Manages deployments to our Kubernetes cluster",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stdout)
			log.SetFormatter(&log.TextFormatter{})

			app.Config = settingsBuilder()
		},
	}

	settingsBuilder = settings.NewBuilder(cmd.PersistentFlags())

	cmd.PersistentFlags().BoolVar(
		&app.IgnoreDeployFailures,
		"ignore-deploy-failures", false,
		"continue deploying services, if any service fails")
	cmd.PersistentFlags().StringSliceVarP(
		&app.Goals,
		"goal", "g", nil,
		"select the goals to execute [all fetch render deploy]")

	cmd.AddCommand(NewBulkCommand(app))
	cmd.AddCommand(NewDeployCommand(app))
	cmd.AddCommand(NewVersionCommand())

	return cmd
}
