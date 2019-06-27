package imagechecker

import "time"

type Options struct {
	WaitTimeout   time.Duration `yaml:"waitTimeout"`
	CheckInterval time.Duration `yaml:"checkInterval"`
	CheckTimeout  time.Duration `yaml:"checkTimeout"`
}

var DefaultOptions = Options{
	WaitTimeout:   10 * time.Minute,
	CheckInterval: 15 * time.Second,
	CheckTimeout:  10 * time.Second,
}
