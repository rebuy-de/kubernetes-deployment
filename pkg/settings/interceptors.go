package settings

type Interceptors struct {
	PreStopSleep        PreStopSleepInterceptor    `yaml:"preStopSleep"`
	RemoveResourceSpecs Interceptor                `yaml:"removeResourceSpecs"`
	Waiter              Interceptor                `yaml:"waiter"`
	GHStatusChecker     GHStatusCheckerInterceptor `yaml:"ghStatusChecker"`
}

type Interceptor struct {
	Enabled TriState `yaml:"enabled"`
}

type PreStopSleepOptions struct {
	Seconds int `yaml:"seconds"`
}

type PreStopSleepInterceptor struct {
	Enabled TriState            `yaml:"enabled"`
	Options PreStopSleepOptions `yaml:"options"`
}

type GHStatusCheckerOptions struct {
	TargetURLRegex string `yaml:"targetUrlRegex"`
	JobRegex       string `yaml:"jobRegex"`
}

type GHStatusCheckerInterceptor struct {
	Enabled TriState               `yaml:"enabled"`
	Options GHStatusCheckerOptions `yaml:"options"`
}

type TriState int

const (
	Unknown TriState = iota
	Disabled
	Enabled
)

func (ts TriState) MarshalYAML() (interface{}, error) {
	switch ts {
	case Enabled:
		return true, nil
	case Disabled:
		return false, nil
	default:
		return nil, nil
	}
}

func (ts *TriState) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var original bool

	err := unmarshal(&original)
	if err != nil {
		return err
	}

	switch original {
	case true:
		(*ts) = Enabled
	case false:
		(*ts) = Disabled
	default:
		(*ts) = Unknown
	}

	return nil
}
