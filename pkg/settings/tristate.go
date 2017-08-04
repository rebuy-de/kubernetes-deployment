package settings

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
