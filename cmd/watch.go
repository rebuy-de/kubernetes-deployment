package cmd

import (
	"fmt"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewWatchCommand(params *api.Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch LABEL_SELECTOR",
		Short: "watches a deployment until it is ready",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("Wrong number of arguments.")
		}

		labelSelector := args[0]

		app, err := api.New(params)
		checkError(err)

		log.WithFields(log.Fields{
			"LabelSelector": labelSelector,
		}).Info("watching project")

		err = app.Watch(labelSelector)
		checkError(err)

		log.WithFields(log.Fields{
			"LabelSelector": labelSelector,
		}).Info("deployment complete")

		return nil
	}

	return cmd
}
