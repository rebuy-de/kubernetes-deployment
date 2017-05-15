package templates

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func ParseManifestFile(inputFile string, outputFile string, settings map[string]string) error {
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	t, err := template.
		New(filepath.Base(inputFile)).
		Funcs(funcMap).
		ParseFiles(inputFile)
	if err != nil {
		return err
	}

	f, err := os.Create(outputFile)
	defer f.Close()
	if err != nil {
		return err
	}

	err = t.Execute(f, settings)

	if err != nil {
		return err
	}

	return nil
}
