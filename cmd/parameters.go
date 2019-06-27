package cmd

import (
	"github.com/aws/aws-sdk-go/aws/session"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/rebuy-de/kubernetes-deployment/pkg/api"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubectl"
	"github.com/rebuy-de/kubernetes-deployment/pkg/settings"
	"github.com/rebuy-de/kubernetes-deployment/pkg/statsdw"
)

const (
	FlagKubeconfig    = "kubeconfig"
	FlagKubectlPath   = "kubectl-path"
	FlagGitHubToken   = "github-token"
	FlagGELFAddress   = "gelf-address"
	FlagStatsdAddress = "statsd-address"
)

const (
	ConfigDir = "~/.rebuy/kubernetes-deployment"
)

type Parameters struct {
	Kubeconfig  string
	KubectlPath string `mapstructure:"kubectl-path"`

	GitHubToken   string `mapstructure:"github-token"`
	GELFAddress   string `mapstructure:"gelf-address"`
	StatsdAddress string `mapstructure:"statsd-address"`
}

func (p *Parameters) Bind(cmd *cobra.Command) {
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
}

func (p *Parameters) ReadIn() error {
	path, err := homedir.Expand(ConfigDir)
	if err != nil {
		return err
	}

	viper.SetConfigName("default")
	viper.AddConfigPath(path)
	viper.ReadInConfig()

	return viper.Unmarshal(p)
}

func (p *Parameters) Kubernetes() (kubernetes.Interface, error) {
	config, err := clientcmd.BuildConfigFromFlags("", p.Kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load kubernetes config")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize kubernetes client")
	}

	return client, nil
}

func (p *Parameters) Build() (*api.App, error) {
	app := new(api.App)

	var err error

	app.Statsd = statsdw.New(p.StatsdAddress)
	app.GitHub = gh.New(p.GitHubToken, app.Statsd)
	app.Kubectl = kubectl.New(p.KubectlPath, p.Kubeconfig)
	app.AWS, err = session.NewSession()
	if err != nil {
		return nil, err
	}

	app.Kubernetes, err = p.Kubernetes()
	if err != nil {
		return nil, err
	}

	app.Settings, err = settings.Read(app.Kubernetes)
	if err != nil {
		return nil, err
	}

	app.Settings.Clean()

	return app, nil
}
