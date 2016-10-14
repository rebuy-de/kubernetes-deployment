package templates

import (
	"os"
	"text/template"
)

func ParseManifestFile(inputFile string, outputFile string, settings map[string]string) error {

	t, err := template.ParseFiles(inputFile)
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
