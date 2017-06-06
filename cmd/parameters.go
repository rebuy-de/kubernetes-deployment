package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagKubeconfig   = "kubeconfig"
	FlagGitHubToken  = "github-token"
	FlagFilename     = "filename"
	FlagHTTPCacheDir = "http-cache-dir"
)

const (
	ConfigDir = "~/.rebuy/kubernetes-deployment"
)

func BindParameters(cmd *cobra.Command) {
	// kubeconfig
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

	// http-cache-dir
	cmd.PersistentFlags().String(
		FlagHTTPCacheDir, "/tmp/kubernetes-deployment-cache",
		"cache directory for HTTP client requests ($HTTP_CACHE_DIR)")
	viper.BindPFlag(FlagHTTPCacheDir, cmd.PersistentFlags().Lookup(FlagHTTPCacheDir))
	viper.BindEnv(FlagHTTPCacheDir, "HTTP_CACHE_DIR")

	// filename
	cmd.PersistentFlags().StringP(
		FlagFilename, "f", "",
		"path to service definitions; might start with './' for local file or 'github.com' for files on GitHub")
	viper.BindPFlag(FlagFilename, cmd.PersistentFlags().Lookup(FlagFilename))
}

func ReadInParameters(p *api.Parameters) error {
	path, err := homedir.Expand(ConfigDir)
	if err != nil {
		return err
	}

	viper.SetConfigName("default")
	viper.AddConfigPath(path)
	viper.ReadInConfig()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.MergeInConfig()

	return viper.Unmarshal(p)
}
