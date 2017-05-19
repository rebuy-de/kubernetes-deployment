package cmd

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes-deployment",
		Short: "Manages deployments to our Kubernetes cluster",
	}

	debug := false
	cmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "more logs")

	params := new(Parameters)
	params.Bind(cmd)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		log.WithFields(log.Fields{
			"version": BuildVersion,
			"date":    BuildDate,
			"commit":  BuildHash,
		}).Infof("kubernetes-deployment started")

		err := params.ReadIn()
		if err != nil {
			log.Fatal(err)
		}

		if strings.TrimSpace(params.Filename) == "" {
			return fmt.Errorf("You have to specify a filename.")
		}

		log.WithFields(
			log.Fields(structs.Map(params)),
		).Debug("config loaded")

		return nil
	}

	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewDeployCommand(params))

	return cmd
}
