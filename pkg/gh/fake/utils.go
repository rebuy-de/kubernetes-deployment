package fake

import (
	"io/ioutil"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

func YAML(obj interface{}) string {
	raw, err := yaml.Marshal(&obj)
	if err != nil {
		panic(err)
	}

	return string(raw)
}

func ScanFiles(root string) Files {
	files := make(Files)

	err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if f.IsDir() {
			return nil
		}

		name, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		raw, err := ioutil.ReadFile(path)
		files[name] = string(raw)
		return err
	})
	if err != nil {
		panic(err)
	}

	return files
}
