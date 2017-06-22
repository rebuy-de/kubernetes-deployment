package cmd

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	log "github.com/sirupsen/logrus"
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

		log.WithFields(log.Fields{
			"Project": project,
			"Branch":  branch,
		}).Info("deploying project")

		err = app.Apply(project, branch)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		log.WithFields(log.Fields{
			"Project": project,
			"Branch":  branch,
		}).Info("updated Kubernetes")

		return nil
	}

	return cmd
}
