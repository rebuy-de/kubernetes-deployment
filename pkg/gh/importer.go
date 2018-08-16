package gh

import (
	"path"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
)

type jsonnetImporter struct {
	client Interface
}

func NewJsonnetImporter(client Interface) jsonnet.Importer {
	return &jsonnetImporter{
		client: client,
	}
}

func (i *jsonnetImporter) Import(importedFrom, importedPath string) (contents jsonnet.Contents, foundAt string, err error) {
	fromLocation, err := NewLocation(importedFrom)
	if err != nil {
		return jsonnet.MakeContents(""), "", err
	}

	var location *Location
	if strings.HasPrefix(importedPath, "github.com/") {
		// location constrution for absolute paths
		location, err = NewLocation(importedPath)
		if err != nil {
			return jsonnet.MakeContents(""), "", err
		}
	} else {
		// location constrution for relative paths
		fromDir := filepath.Dir(fromLocation.String())
		absolute := path.Clean(path.Join(fromDir, importedPath))

		location, err = NewLocation(absolute)
		if err != nil {
			return jsonnet.MakeContents(""), "", err
		}
	}

	// having no ref means same ref, if inside same repo and "master" if it is in another repo
	if location.Ref == "" {
		if location.Owner == fromLocation.Owner && location.Repo == fromLocation.Repo {
			location.Ref = fromLocation.Ref
		} else {
			location.Ref = "master"
		}

	}

	file, err := i.client.GetFile(location)
	if err != nil {
		return jsonnet.MakeContents(""), "", err
	}

	return jsonnet.MakeContents(file.Content), location.String(), nil
}
