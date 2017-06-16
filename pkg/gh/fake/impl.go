package fake

import (
	"path"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

func (d *GitHub) GetBranch(l *gh.Location) (*gh.Branch, error) {
	branch := (*d)[l.Owner][l.Repo][l.Ref].Meta
	return &branch, nil
}

func (d *GitHub) GetFile(l *gh.Location) (string, error) {
	content := (*d)[l.Owner][l.Repo][l.Ref].Files[l.Path]
	return content, nil
}

func (d *GitHub) GetFiles(l *gh.Location) (map[string]string, error) {
	files := make(map[string]string)

	for p, content := range (*d)[l.Owner][l.Repo][l.Ref].Files {
		dir, file := path.Split("/" + p)
		if path.Clean("/"+dir+"/") == path.Clean("/"+l.Path+"/") {
			files[path.Clean(file)] = content
		}
	}

	return files, nil
}
