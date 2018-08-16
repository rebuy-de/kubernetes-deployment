package fake

import (
	"fmt"
	"path"

	"github.com/google/go-github/github"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

func (d *GitHub) repos(l *gh.Location) (Repos, error) {
	repos, ok := (*d)[l.Owner]
	if !ok {
		return nil, fmt.Errorf("fake owner '%s' doesn't exist", l.Owner)
	}

	return repos, nil
}

func (d *GitHub) branches(l *gh.Location) (Branches, error) {
	repos, err := d.repos(l)
	if err != nil {
		return nil, err
	}

	branches, ok := repos[l.Repo]
	if !ok {
		return nil, fmt.Errorf("fake repo '%s/%s' doesn't exist", l.Owner, l.Repo)
	}

	return branches, nil
}

func (d *GitHub) ref(l *gh.Location) (Branch, error) {
	branches, err := d.branches(l)
	if err != nil {
		return Branch{}, err
	}

	ref, ok := branches[l.Ref]
	if !ok {
		return Branch{}, fmt.Errorf("fake branch '%s/%s#%s' doesn't exist", l.Owner, l.Repo, l.Ref)
	}

	return ref, nil
}

func (d *GitHub) GetBranch(l *gh.Location) (*gh.Branch, error) {
	ref, err := d.ref(l)
	if err != nil {
		return nil, err
	}

	branch := ref.Meta
	return &branch, nil
}

func (d *GitHub) GetFile(l *gh.Location) (gh.File, error) {
	for _, file := range (*d)[l.Owner][l.Repo][l.Ref].Files {
		if file.Location.Path == l.Path {
			return file, nil
		}
	}
	return gh.File{}, fmt.Errorf("File %s not found", l.String())
}

func (d *GitHub) GetFiles(l *gh.Location) ([]gh.File, error) {
	var files []gh.File

	for _, file := range (*d)[l.Owner][l.Repo][l.Ref].Files {
		dir, _ := path.Split("/" + file.Location.Path)
		if path.Clean("/"+dir+"/") == path.Clean("/"+l.Path+"/") {
			files = append(files, file)
		}
	}

	return files, nil
}

func (d *GitHub) GetStatuses(location *gh.Location) ([]github.RepoStatus, error) {
	return nil, fmt.Errorf("not implemented, yet")
}
