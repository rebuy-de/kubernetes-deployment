package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewGenerateCommand(params *api.Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate PROJECT [BRANCH]",
		Short: "Views the generated manifests for a project",
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
		}).Info("generating manifests")

		objects, err := app.Generate(project, branch)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		for _, obj := range objects {
			raw, err := json.MarshalIndent(obj, "", "    ")
			if err != nil {
				log.Fatal(err)
				return nil
			}
			fmt.Println(string(raw))
		}

		return nil
	}

	return cmd
}
