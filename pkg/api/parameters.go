package api

type Parameters struct {
	Context string

	Kubeconfig  string
	KubectlPath string `mapstructure:"kubectl-path"`

	GitHubToken  string `mapstructure:"github-token"`
	HTTPCacheDir string `mapstructure:"http-cache-dir"`
	GELFAddress  string `mapstructure:"gelf-address"`

	Filename string
}
