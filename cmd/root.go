package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes-deployment",
		Short: "Manages deployments to our Kubernetes cluster",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stdout)
		},
	}

	params := new(Parameters)
	params.Bind(cmd)

	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewDeployCommand(params))

	return cmd
}
