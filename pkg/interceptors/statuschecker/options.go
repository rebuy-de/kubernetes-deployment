package statuschecker

type Options struct {
	TargetURLRegex string `yaml:"targetUrlRegex"`
	JobRegex       string `yaml:"jobRegex"`
}
