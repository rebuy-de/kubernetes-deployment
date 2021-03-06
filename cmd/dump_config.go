package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDumpConfigCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dump-config",
		Short: "Dumps the current configuration",
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		all := viper.AllSettings()

		raw, err := json.MarshalIndent(all, "", "    ")
		cmdutil.Must(err)
		fmt.Println(string(raw))
	}

	return cmd
}
