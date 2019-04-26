package cmd

import (
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewApplyCommand(params *Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply PROJECT [BRANCH]",
		Short: "Deploys a project to Kubernetes",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		project, branch, err := getProject(args)
		if err != nil {
			return err
		}

		app, err := params.Build()
		cmdutil.Must(err)
		defer app.Close()

		log.WithFields(log.Fields{
			"Project": project,
			"Branch":  branch,
		}).Info("deploying project")

		err = app.Apply(project, branch)
		cmdutil.Must(err)

		log.WithFields(log.Fields{
			"Project": project,
			"Branch":  branch,
		}).Info("deployment finished")

		return nil
	}

	return cmd
}
