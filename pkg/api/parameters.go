package api

type Parameters struct {
	Kubeconfig  string
	KubectlPath string `mapstructure:"kubectl-path"`

	GitHubToken   string `mapstructure:"github-token"`
	HTTPCacheDir  string `mapstructure:"http-cache-dir"`
	GELFAddress   string `mapstructure:"gelf-address"`
	StatsdAddress string `mapstructure:"statsd-address"`

	Filename string
}
