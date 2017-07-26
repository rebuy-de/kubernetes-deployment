package testutil

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var (
	UpdateGolden = flag.Bool("update-golden", false,
		"update the golden file instead of comparing it")
)

func AssertGoldenYAML(t *testing.T, filename string, data interface{}) {
	generated, err := yaml.Marshal(data)
	if err != nil {
		t.Error(err)
		return
	}

	generated = append(generated, '\n')

	AssertGolden(t, filename, generated)
}

func AssertGoldenJSON(t *testing.T, filename string, data interface{}) {
	generated, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Error(err)
		return
	}

	generated = append(generated, '\n')

	AssertGolden(t, filename, generated)
}

func AssertGolden(t *testing.T, filename string, data []byte) {
	if *UpdateGolden {
		err := ioutil.WriteFile(filename, data, os.FileMode(0644))
		if err != nil {
			t.Error(err)
			return
		}
	}

	golden, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Error(err)
		return
	}

	if string(golden) != string(data) {
		t.Errorf("Generated file '%s' doesn't match golden file. Update with '-update-golden'.", filename)
	}
}
