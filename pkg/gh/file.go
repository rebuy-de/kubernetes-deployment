package gh

import "path"

type File struct {
	Location *Location
	Content  string
}

func (f *File) Name() string {
	_, name := path.Split(f.Location.Path)
	return name
}
