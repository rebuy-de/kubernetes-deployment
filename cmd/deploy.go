package cmd

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
)

const (
	DefaultBranch = "master"
)

func NewDeployCommand(params *Parameters) *cobra.Command {
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

		settings := params.LoadSettings()

		log.Infof("deploying %s/%s", project, branch)

		service := settings.Service(project)
		if service == nil {
			log.WithFields(log.Fields{
				"Project": project,
			}).Fatal("project not found")
		}
		log.WithFields(
			log.Fields(structs.Map(service)),
		).Debug("service found")

		templates, err := params.GitHubClient().GetFiles(service.Location)
		if err != nil {
			log.Error(err)
			return nil
		}

		fmt.Printf("%#v\n", templates)

		return nil
	}

	return cmd
}
