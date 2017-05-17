package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagKubeconfig  = "kubeconfig"
	FlagGitHubToken = "github-token"
	FlagFilename    = "filename"
)

type Parameters struct {
	Kubeconfig  string
	GitHubToken string
	Filename    string
}

func (p *Parameters) Bind(cmd *cobra.Command) {
	// kubeconfig
	viper.SetDefault(FlagKubeconfig, "$HOME/.kube/config")
	cmd.PersistentFlags().String(
		FlagKubeconfig, "",
		"path to the kubeconfig file to use for deployments ($KUBECONFIG)")
	viper.BindPFlag(FlagKubeconfig, cmd.PersistentFlags().Lookup(FlagKubeconfig))
	viper.BindEnv(FlagKubeconfig, "KUBECONFIG")

	// github-token
	cmd.PersistentFlags().String(
		FlagGitHubToken, "",
		"oauth token for GitHub ($GITHUB_TOKEN)")
	viper.BindPFlag(FlagGitHubToken, cmd.PersistentFlags().Lookup(FlagGitHubToken))
	viper.BindEnv(FlagGitHubToken, "GITHUB_TOKEN")

	// filename
	cmd.PersistentFlags().StringP(
		FlagFilename, "f", "",
		"path to service definitions; might start with './' for local file or 'github.com' for files on GitHub")
	viper.BindPFlag(FlagFilename, cmd.PersistentFlags().Lookup(FlagFilename))
}

func (p *Parameters) ReadIn() error {
	return viper.Unmarshal(p)
}
