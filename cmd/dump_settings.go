package cmd

import (
	"fmt"

	yaml "gopkg.in/yaml.v2"

	log "github.com/Sirupsen/logrus"
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
		if err != nil {
			log.Fatal(err)
			return
		}

		app.Settings.Clean(app.Parameters.Context)

		raw, err := yaml.Marshal(app.Settings)
		if err != nil {
			log.Fatal(err)
			return
		}
		fmt.Println(string(raw))
	}

	return cmd
}
