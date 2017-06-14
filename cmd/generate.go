package cmd

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
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

		objects, err := api.Generate(params, project, branch)
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