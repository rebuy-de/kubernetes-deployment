package cmd

import (
	log "github.com/Sirupsen/logrus"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	FlagKubeconfig  = "kubeconfig"
	FlagGitHubToken = "github-token"
	FlagFilename    = "filename"
)

const (
	ConfigDir = "~/.rebuy/kubernetes-deployment"
)

type Parameters struct {
	Kubeconfig  string
	GitHubToken string `mapstructure:"github-token"`
	Filename    string
}

func (p *Parameters) Bind(cmd *cobra.Command) {
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

	// filename
	cmd.PersistentFlags().StringP(
		FlagFilename, "f", "",
		"path to service definitions; might start with './' for local file or 'github.com' for files on GitHub")
	viper.BindPFlag(FlagFilename, cmd.PersistentFlags().Lookup(FlagFilename))
}

func (p *Parameters) ReadIn() error {
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

func (p *Parameters) LoadSettings() *settings.Settings {
	ghClient := gh.New(p.GitHubToken)

	sett, err := settings.Read(p.Filename, ghClient)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"ServiceCount": len(sett.Services),
		"Filename":     p.Filename,
	}).Debug("loaded service file")

	return sett
}
