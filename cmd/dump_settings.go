package cmd

import (
	"fmt"

	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

func NewDumpSettingsCommand(params *Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump-settings",
		Short: "Dumps the current settings",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		app, err := params.Build()
		cmdutil.Must(err)
		defer app.Close()

		raw, err := yaml.Marshal(app.Settings)
		cmdutil.Must(err)
		fmt.Println(string(raw))
	}

	return cmd
}
