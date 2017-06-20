package fake

import "github.com/rebuy-de/kubernetes-deployment/pkg/gh"

type GitHub map[string]Repos

type Repos map[string]Branches

type Branches map[string]Branch

type Branch struct {
	Meta  gh.Branch
	Files Files
}

type Files map[string]string
