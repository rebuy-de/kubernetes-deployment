package settings

type Services []*Service

func (s *Services) Clean() {
	for _, service := range *s {
		service.Clean()
	}
}
