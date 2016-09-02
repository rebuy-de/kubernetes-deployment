package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

const (
	DEFAULT_BRANCH = "master"
	DEFAULT_PATH   = "deployment/kubernetes"
)

type Service struct {
	Repository string `yaml:"repo"`
	Path       string `yaml:"path,omitempty"`
	Branch     string `yaml:"branch,omitempty"`
}

func (s *Service) Clean() {
	if strings.TrimSpace(s.Path) == "" {
		s.Path = DEFAULT_PATH
	}

	if strings.TrimSpace(s.Branch) == "" {
		s.Branch = DEFAULT_BRANCH
	}

	if strings.HasPrefix(s.Repository, "github.com/") {
		base := strings.TrimPrefix(s.Repository, "github.com/")
		s.Repository = fmt.Sprintf("git@github.com:%s.git", base)
	}

	s.Path = filepath.Clean("/" + s.Path)
}
