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
	Name       string `yaml:"name,omitempty"`
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
	s.Name = s.GuessName()
}

func (s *Service) GuessName() string {
	parts := []string{}

	repo := s.Repository
	repo = strings.TrimPrefix(repo, "git@github.com:rebuy-de/")
	repo = strings.TrimSuffix(repo, ".git")
	repo = strings.Replace(repo, "/", "-", -1)
	parts = append(parts, repo)

	if s.Path != DEFAULT_PATH {
		path := s.Path
		path = strings.Trim(path, "/")
		path = strings.Replace(path, "/", "-", -1)
		parts = append(parts, path)
	}

	if s.Branch != DEFAULT_BRANCH {
		parts = append(parts, s.Branch)
	}

	return strings.Join(parts, "-")
}
