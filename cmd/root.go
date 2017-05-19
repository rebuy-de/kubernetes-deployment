package cmd

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes-deployment",
		Short: "Manages deployments to our Kubernetes cluster",
	}

	params := new(Parameters)
	params.Bind(cmd)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.DebugLevel)
		log.SetOutput(os.Stdout)

		err := params.ReadIn()
		if err != nil {
			log.Fatal(err)
		}

		if strings.TrimSpace(params.Filename) == "" {
			return fmt.Errorf("You have to specify a filename.")
		}

		return nil
	}

	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewDeployCommand(params))

	return cmd
}
