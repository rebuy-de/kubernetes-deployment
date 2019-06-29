package cmd

import (
	"fmt"
	"strings"

	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"

	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	cmdutilv2 "github.com/rebuy-de/rebuy-go-sdk/v2/cmdutil"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := cmdutilv2.New(
		"kubernetes-deployment", "Manages deployments to our Kubernetes cluster",
	)

	debug := false
	cmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "show debug log messages")

	jsonLogs := false
	cmd.PersistentFlags().BoolVar(&jsonLogs, "json-logs", false, "prints the logs as JSON")

	params := new(Parameters)
	params.Bind(cmd)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		if jsonLogs {
			log.SetFormatter(&log.JSONFormatter{
				FieldMap: log.FieldMap{
					log.FieldKeyTime:  "Time",
					log.FieldKeyLevel: "Level",
					log.FieldKeyMsg:   "Message",
				},
			})
		}

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

		log.WithFields(log.Fields{
			"Version": cmdutil.BuildVersion,
			"Date":    cmdutil.BuildDate,
			"Commit":  cmdutil.BuildHash,
		}).Debugf("kubernetes-deployment started")

		if strings.TrimSpace(params.GitHubToken) == "" {
			return fmt.Errorf("You have to specify a GitHubToken.")
		}

		log.WithFields(log.Fields{
			"GitHubToken": fmt.Sprintf("%s****", params.GitHubToken[0:4]),
			"Kubeconfig":  params.Kubeconfig,
		}).Debug("config loaded")

		return nil
	}

	cmd.AddCommand(cmdutil.NewVersionCommand())
	cmd.AddCommand(NewDumpConfigCommand())
	cmd.AddCommand(NewDumpSettingsCommand(params))
	cmd.AddCommand(NewApplyCommand(params))
	cmd.AddCommand(NewServerCommand(params))

	return cmd
}
