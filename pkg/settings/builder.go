package settings

import (
	"fmt"

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
	stage := fs.StringP(
		"stage", "s", "testing",
		"name of the stage; used for selecting the configuration (eg ~/.rebuy/kubernetes-deployment/<stage>.yaml)")

	v := viper.New()

	v.SetDefault(kubeconfig, "$HOME/.kube/config")
	v.SetDefault(output, "./output")

	fs.String(kubeconfigFlag, "", "path to the kubeconfig file to use for deployments ($KUBECONFIG)")
	v.BindPFlag(kubeconfig, fs.Lookup(kubeconfigFlag))
	v.BindEnv(kubeconfig, "KUBECONFIG")

	return func() ProjectConfig {
		if localConfigPath != nil {
			v.SetConfigFile(*localConfigPath)
			v.ReadInConfig()
		}

		v.SetConfigName("config")
		v.AddConfigPath(fmt.Sprintf("$HOME/.rebuy/kubernetes-deployment/%s.yaml", *stage))
		v.MergeInConfig()

		s := ProjectConfig{}
		v.Unmarshal(&s)
		return s
	}
}
