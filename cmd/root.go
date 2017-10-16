package cmd

import (
	"fmt"
	"strings"

	graylog "gopkg.in/gemnasium/logrus-graylog-hook.v2"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes-deployment",
		Short: "Manages deployments to our Kubernetes cluster",
	}

	debug := false
	cmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "show debug log messages")

	jsonLogs := false
	cmd.PersistentFlags().BoolVar(&jsonLogs, "json-logs", false, "prints the logs as JSON")

	params := BindParameters(cmd)

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

		err := ReadInParameters(params)
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
			"Version": BuildVersion,
			"Date":    BuildDate,
			"Commit":  BuildHash,
		}).Debugf("kubernetes-deployment started")

		if strings.TrimSpace(params.GitHubToken) == "" {
			return fmt.Errorf("You have to specify a GitHubToken.")
		}

		if strings.TrimSpace(params.Kubeconfig) == "" {
			return fmt.Errorf("You have to specify a path to kubeconfig.")
		}

		if strings.TrimSpace(params.Filename) == "" {
			return fmt.Errorf("You have to specify a filename.")
		}

		log.WithFields(log.Fields{
			"GitHubToken": fmt.Sprintf("%s****", params.GitHubToken[0:4]),
			"Kubeconfig":  params.Kubeconfig,
			"Filename":    params.Filename,
		}).Debug("config loaded")

		return nil
	}

	cmd.AddCommand(NewVersionCommand())
	cmd.AddCommand(NewDumpConfigCommand())
	cmd.AddCommand(NewDumpSettingsCommand(params))
	cmd.AddCommand(NewApplyCommand(params))
	cmd.AddCommand(NewGenerateCommand(params))

	return cmd
}
