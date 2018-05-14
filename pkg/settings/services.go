package settings

type Services []Service

func (s Services) Get(name string) *Service {
	for _, service := range s {
		if service.Name == name {
			return &service
		}
	}

	// second loop, to prioritise names
	for _, service := range s {
		for _, alias := range service.Aliases {
			if alias == name {
				return &service
			}
		}
	}

	return nil
}
