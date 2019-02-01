package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagKubeconfig    = "kubeconfig"
	FlagKubectlPath   = "kubectl-path"
	FlagGitHubToken   = "github-token"
	FlagFilename      = "filename"
	FlagHTTPCacheDir  = "http-cache-dir"
	FlagGELFAddress   = "gelf-address"
	FlagStatsdAddress = "statsd-address"
)

const (
	ConfigDir = "~/.rebuy/kubernetes-deployment"
)

func BindParameters(cmd *cobra.Command) *api.Parameters {
	params := new(api.Parameters)

	// kubeconfig
	cmd.PersistentFlags().String(
		FlagKubeconfig, "",
		"path to the kubeconfig file to use for deployments ($KUBECONFIG)")
	viper.BindPFlag(FlagKubeconfig, cmd.PersistentFlags().Lookup(FlagKubeconfig))
	viper.BindEnv(FlagKubeconfig, "KUBECONFIG")

	// kubectl-path
	cmd.PersistentFlags().String(
		FlagKubectlPath, "kubectl",
		"full path or basename of kubectl executable")
	viper.BindPFlag(FlagKubectlPath, cmd.PersistentFlags().Lookup(FlagKubectlPath))

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

	// gelf address
	cmd.PersistentFlags().String(
		FlagGELFAddress, "",
		"a Graylog GELF UDP address (a ip:port string) for sending logs")
	viper.BindPFlag(FlagGELFAddress, cmd.PersistentFlags().Lookup(FlagGELFAddress))

	// statsd address
	cmd.PersistentFlags().String(
		FlagStatsdAddress, "",
		"a statsd UDP address (a ip:port string) for sending metrics")
	viper.BindPFlag(FlagStatsdAddress, cmd.PersistentFlags().Lookup(FlagStatsdAddress))

	return params
}

func ReadInParameters(p *api.Parameters) error {
	path, err := homedir.Expand(ConfigDir)
	if err != nil {
		return err
	}

	viper.SetConfigName("default")
	viper.AddConfigPath(path)
	viper.ReadInConfig()

	return viper.Unmarshal(p)
}
