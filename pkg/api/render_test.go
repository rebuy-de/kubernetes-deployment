package api

import (
	"fmt"
	"io/ioutil"
	"path"
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	fakeGH "github.com/rebuy-de/kubernetes-deployment/pkg/gh/fake"
	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"
	"github.com/rebuy-de/rebuy-go-sdk/testutil"
)

func readFile(t *testing.T, path string) string {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(dat)
}

func TestDecode(t *testing.T) {
	cases := []struct {
		name  string
		files []string
	}{
		{
			name: "multi-yaml",
			files: []string{
				"render-deployment.yaml",
				"render-podpreset.yaml",
			},
		},
		{
			name: "simple-jsonnet",
			files: []string{
				"render-pdb.jsonnet",
			},
		},
		{
			name: "complex-jsonnet",
			files: []string{
				"render-job.jsonnet",
			},
		},
		{
			name: "local-import-jsonnet",
			files: []string{
				"render-silo.jsonnet",
				"render-silo.libsonnet",
			},
		},
		{
			name: "unknown-object-yaml",
			files: []string{
				"unknown-object.yaml",
			},
		},
		{
			name: "unknown-object-jsonnet",
			files: []string{
				"unknown-object.jsonnet",
			},
		},
	}

	vars := templates.Variables{
		"testString":  "bish-bash-bosh",
		"gitCommitID": "ffffff",
		"clusterName": "staging",
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			files := fakeGH.Files{}
			for _, fname := range tc.files {
				files = append(files, gh.File{
					Location: &gh.Location{
						Owner: "rebuy-de",
						Repo:  "test",
						Path:  fname,
						Ref:   "master",
					},
					Content: readFile(t, path.Join("test-fixtures", fname)),
				})
			}

			github := &fakeGH.GitHub{
				"rebuy-de": fakeGH.Repos{
					"test": fakeGH.Branches{
						"master": fakeGH.Branch{
							Meta: gh.Branch{
								Name: "master",
								SHA:  "5a5369823a2a9a6ad9c241b404be39f802d41d48",
							},
							Files: files,
						},
					},
				},
			}

			app := App{
				Interceptors: &interceptors.Multi{},
				Clients: &Clients{
					GitHub: github,
				},
			}

			objects, err := app.decode(files, vars)
			if err != nil {
				t.Fatalf("%+v", err)
			}
			g := path.Join("test-fixtures", fmt.Sprintf("render-golden-%s.json", tc.name))
			testutil.AssertGoldenJSON(t, g, objects)
		})
	}
}
