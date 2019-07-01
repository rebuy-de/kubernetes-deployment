package cmd

import (
	"fmt"
	"strings"

	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"

	cmdutilv2 "github.com/rebuy-de/rebuy-go-sdk/v2/cmdutil"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func WithLogJSON() cmdutilv2.Option {
	var (
		enabled bool
	)

	return func(cmd *cobra.Command) error {
		cmd.PersistentFlags().BoolVar(
			&enabled, "json-logs", false, "prints the logs as JSON")

		cmd.PreRun = func(cmd *cobra.Command, args []string) {
			if enabled {
				log.SetFormatter(&log.JSONFormatter{
					FieldMap: log.FieldMap{
						log.FieldKeyTime:  "Time",
						log.FieldKeyLevel: "Level",
						log.FieldKeyMsg:   "Message",
					},
				})
			}
		}

		return nil
	}
}

func NewRootCommand() *cobra.Command {
	cmd := cmdutilv2.New(
		"kubernetes-deployment", "Manages deployments to our Kubernetes cluster",
		cmdutilv2.WithVersionCommand(),
		cmdutilv2.WithVersionLog(logrus.DebugLevel),
		cmdutilv2.WithLogVerboseFlag(),
		WithLogJSON(),
	)

	params := new(Parameters)
	params.Bind(cmd)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {

		err := params.ReadIn()
		if err != nil {
			log.Fatal(err)
		}

		if params.GELFAddress != "" {
			hook := graylog.NewGraylogHook(params.GELFAddress,
				map[string]interface{}{
					"run":      randomID(),
					"facility": "kubernetes-deployment",
				})
			hook.Level = log.DebugLevel
			log.AddHook(hook)
		}

		if strings.TrimSpace(params.GitHubToken) == "" {
			return fmt.Errorf("You have to specify a GitHubToken.")
		}

		log.WithFields(log.Fields{
			"GitHubToken": fmt.Sprintf("%s****", params.GitHubToken[0:4]),
			"Kubeconfig":  params.Kubeconfig,
		}).Debug("config loaded")

		return nil
	}

	cmd.AddCommand(NewDumpConfigCommand())
	cmd.AddCommand(NewDumpSettingsCommand(params))
	cmd.AddCommand(NewApplyCommand(params))
	cmd.AddCommand(NewServerCommand(params))

	return cmd
}
