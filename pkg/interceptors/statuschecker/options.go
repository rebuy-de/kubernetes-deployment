package statuschecker

import "time"

type Options struct {
	TargetURLRegex string `yaml:"targetUrlRegex"`
	ContextRegex   string `yaml:"contextRegex"`

	InitialDelay         time.Duration `yaml:"initialDelay"`
	PullInterval         time.Duration `yaml:"pullInterval"`
	NotificationInterval time.Duration `yaml:"notificationInterval"`
}

var DefaultOptions = Options{
	TargetURLRegex: `.*`,
	ContextRegex:   `.*`,

	InitialDelay:         20 * time.Second,
	PullInterval:         10 * time.Second,
	NotificationInterval: 3 * time.Minute,
}
