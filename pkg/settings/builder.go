package settings

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type SettingsBuilder func() ProjectConfig

func NewBuilder(fs *pflag.FlagSet) SettingsBuilder {
	var (
		kubeconfigFlag = "kubeconfig"
		kubeconfig     = "settings.kubeconfig"
		output         = "settings.output"
	)

	localConfigPath := fs.StringP(
		"config", "c", "./config/services.yaml",
		"project configuration file")

	v := viper.New()

	v.SetDefault(kubeconfig, "$HOME/.kube/config")
	v.SetDefault(output, "./output")

	fs.String(kubeconfigFlag, "", "path to the kubeconfig file to use for deployments ($KUBECONFIG)")
	v.BindPFlag(kubeconfig, fs.Lookup(kubeconfigFlag))
	v.BindEnv(kubeconfig, "KUBECONFIG")

	v.SetConfigName("config")
	v.AddConfigPath("$HOME/.kubernetes-deployment")
	v.ReadInConfig()

	return func() ProjectConfig {
		if localConfigPath != nil {
			v.SetConfigFile(*localConfigPath)
			v.MergeInConfig()
		}

		s := ProjectConfig{}
		v.Unmarshal(&s)
		return s
	}
}
