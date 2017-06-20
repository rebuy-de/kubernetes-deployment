package testutil

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

var (
	UpdateGolden = flag.Bool("update-golden", false,
		"update the golden file instead of comparing it")
)

func AssertGoldenFile(t *testing.T, filename string, data interface{}) {
	generated, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		t.Error(err)
		return
	}

	generated = append(generated, '\n')

	if *UpdateGolden {
		err := ioutil.WriteFile(filename, generated, os.FileMode(0644))
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

	if string(golden) != string(generated) {
		t.Errorf("Generated file '%s' doesn't match golden file. Update with '-update-golden'.", filename)
	}
}
