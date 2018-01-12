package kubectl

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"testing"
)

func readFile(t *testing.T, path string) []byte {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return dat
}

func jsonToMap(t *testing.T, b []byte) map[string]interface{} {
	obj := make(map[string]interface{})
	err := json.Unmarshal(b, &obj)
	if err != nil {
		t.Fatal(err)
	}
	return obj
}

func TestEmptyStatus(t *testing.T) {
	dir := "test-fixtures"
	tcs := []string{
		path.Join(dir, "pod-disruption-budget.json"),
		path.Join(dir, "storage-class.json"),
	}

	for _, tc := range tcs {
		t.Run(tc, func(t *testing.T) {
			raw := readFile(t, tc)
			before := jsonToMap(t, raw)
			_, hasStatus := before["status"]

			stripped, err := emptyStatus(raw)
			if err != nil {
				t.Fatal(err)
			}

			m := jsonToMap(t, stripped)
			if m["status"] != nil {
				t.Fatal("JSON still contains status field")
			}
			if !hasStatus && len(before) != len(m) {
				t.Fatal("dropped non status field")
			}
		})
	}
}
