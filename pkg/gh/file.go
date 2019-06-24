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

type FileByName []File

func (a FileByName) Len() int           { return len(a) }
func (a FileByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a FileByName) Less(i, j int) bool { return a[i].Location.String() < a[j].Location.String() }
