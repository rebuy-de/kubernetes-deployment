package api

import (
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/testutil"
	"io/ioutil"
	"path"
	"testing"
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
	tcs := []gh.File{
		{
			Path:    "manifest-deployment.yaml",
			Content: readFile(t, path.Join(dir, "manifest-deployment.yaml")),
		},
		{
			Path:    "manifest-podpreset.yaml",
			Content: readFile(t, path.Join(dir, "manifest-podpreset.yaml")),
		},
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
