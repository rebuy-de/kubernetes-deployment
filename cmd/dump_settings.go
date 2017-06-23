package cmd

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
)

func NewDumpSettingsCommand(params *api.Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump-settings",
		Short: "Dumps the current settings",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		app, err := api.New(params)
		checkError(err)

		app.Settings.Clean(app.Parameters.Context)

		raw, err := yaml.Marshal(app.Settings)
		checkError(err)
		fmt.Println(string(raw))
	}

	return cmd
}
