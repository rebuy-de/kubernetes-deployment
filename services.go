package main

import "math/rand"

type Services []*Service

func (s *Services) Clean() {
	for _, service := range *s {
		service.Clean()
	}
}

func (s *Services) Shuffle() {
	slice := *s
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
