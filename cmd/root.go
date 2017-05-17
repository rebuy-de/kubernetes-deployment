package cmd

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "golang-template",
		Short: "an example app for golang which can be used as template",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			log.SetLevel(log.DebugLevel)
			log.SetOutput(os.Stdout)
		},
	}

	cmd.AddCommand(NewVersionCommand())

	return cmd
}
