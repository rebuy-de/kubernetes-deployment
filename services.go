package main

import (
	"math/rand"
	"time"
)

type Services []*Service

func (s *Services) Clean() {
	for _, service := range *s {
		service.Clean()
	}
}

func (s *Services) Shuffle() {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)

	slice := *s
	for i := range slice {
		j := r.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}
