package cmd

import (
	"context"
	"net/http"

	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewServerCommand(params *Parameters) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "k26r (kubernetes-deployment-server)",
	}

	var httpListenAddress string
	cmd.PersistentFlags().StringVar(&httpListenAddress, "listen-addr", ":8080",
		"listening address for the HTTP server")

	cmd.RunE = func(cmd *cobra.Command, args []string) error {
		app, err := params.Build()
		cmdutil.Must(err)
		defer app.Close()

		controller := new(Controller)
		controller.App = app

		server := &http.Server{
			Addr:    httpListenAddress,
			Handler: controller.Mux(),
		}

		ctx := context.TODO()
		ctx, cancel := context.WithCancel(ctx)

		go func() {
			logrus.Error(server.ListenAndServe())
			cancel()
		}()
		logrus.Info("started HTTP server")

		<-ctx.Done()
		logrus.Warn("shutting down HTTP server")
		logrus.Error(server.Shutdown(context.Background()))

		return nil
	}

	return cmd
}
