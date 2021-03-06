package kubeutil

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/rebuy-de/rebuy-go-sdk/testutil"
	"k8s.io/apimachinery/pkg/runtime"
)

func readFile(t *testing.T, filename string) []byte {
	dat, err := ioutil.ReadFile(path.Join("test-fixtures", filename))
	if err != nil {
		t.Fatal(err)
	}
	return dat
}

func TestUnknownObject_interface(t *testing.T) {
	uobj := new(UnknownObject)

	_ = runtime.Object(uobj)
}

func TestUnknownObject_parseJSON(t *testing.T) {
	raw := readFile(t, "unknown.json")

	uo := new(UnknownObject)
	err := uo.FromJSON(raw)
	if err != nil {
		t.Fatal(err)
	}

	uo.ObjectMeta.Labels["foo"] = "bar"

	testutil.AssertGoldenJSON(t, "test-fixtures/unknown-json-golden.json", uo)
}

func TestUnknownObject_parseYAML(t *testing.T) {
	raw := readFile(t, "unknown.yaml")

	uo := new(UnknownObject)
	err := uo.FromYAML(raw)
	if err != nil {
		t.Fatal(err)
	}

	uo.ObjectMeta.Labels["foo"] = "bar"

	testutil.AssertGoldenJSON(t, "test-fixtures/unknown-yaml-golden.json", uo)
}
