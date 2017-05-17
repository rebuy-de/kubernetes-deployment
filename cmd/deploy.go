package cmd

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewDeployCommand(params *Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy PROJECT [BRANCH]",
		Short: "Deploys a project to Kubernetes",
	}

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := params.ReadIn()
		if err != nil {
			log.Fatal(err)
		}

		if strings.TrimSpace(params.Filename) == "" {
			return fmt.Errorf("You have to specify a filename.")
		}

		if len(args) < 1 || len(args) > 2 {
			return fmt.Errorf("Wrong number of arguments.")
		}

		fmt.Printf("%#v", params)

		return nil
	}

	return cmd
}
