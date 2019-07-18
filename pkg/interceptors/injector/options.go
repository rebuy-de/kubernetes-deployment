package injector

type Options struct {
	InjectArguments []string `yaml:"injectArguments"`
	ConnectTimeout  string   `yaml:"connectTimeout"`
}

var DefaultOptions = Options{
	InjectArguments: []string{"--manual", "--proxy-memory-request", "20Mi", "--proxy-cpu-request", "35m"},
	ConnectTimeout:  "10s",
}
