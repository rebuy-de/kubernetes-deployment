package main

import (
	"github.com/rebuy-de/rebuy-go-sdk/cmdutil"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	defer cmdutil.HandleExit()
	if err := NewRootCommand().Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func NewRootCommand() *cobra.Command {
	cmd := cmdutil.NewRootCommand(new(App))
	cmd.Short = "k26r (kubernetes-deployment-server)"
	return cmd
}
