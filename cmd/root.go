package cmd

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubernetes-deployment",
		Short: "Manages deployments to our Kubernetes cluster",
	}

	debug := false
	cmd.PersistentFlags().BoolVarP(&debug, "verbose", "v", false, "more logs")

	BindParameters(cmd)
	params := new(api.Parameters)

	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		log.SetLevel(log.InfoLevel)
		if debug {
			log.SetLevel(log.DebugLevel)
		}

		log.WithFields(log.Fields{
			"Version": BuildVersion,
			"Date":    BuildDate,
			"Commit":  BuildHash,
		}).Infof("kubernetes-deployment started")

		err := ReadInParameters(params)
		if err != nil {
			log.Fatal(err)
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
	cmd.AddCommand(NewDeployCommand(params))
	cmd.AddCommand(NewGenerateCommand(params))

	return cmd
}
