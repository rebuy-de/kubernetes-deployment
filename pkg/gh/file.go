package gh

import "path"

type File struct {
	Path    string
	Content string
}

func (f *File) Name() string {
	_, name := path.Split(f.Path)
	return name
}
