package api

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/testutil"
	"testing"
	"io/ioutil"
	"path"
)

func readFile(t *testing.T, path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(dat)
}

func TestDecode(t *testing.T) {
	dir := "test-fixtures"
	tcs := map[string]string{
		"manifest-deployment.yaml": readFile(t, path.Join(dir, "manifest-deployment.yaml")),
		"manifest-podpreset.yaml": readFile(t, path.Join(dir, "manifest-podpreset.yaml")),
	}
	golden := "decoded-golden.yaml"

	app := App{Interceptors: &interceptors.Multi{}}
	objects, err := app.decode(tcs)
	if err != nil {
		t.Fatal(err)
	}
	g := path.Join(dir, golden)
	testutil.AssertGoldenJSON(t, g, objects)
}
