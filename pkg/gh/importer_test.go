package gh_test

import (
	"testing"
	"time"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh/fake"
)

func TestImporter(t *testing.T) {
	client := &fake.GitHub{
		"rebuy-de": fake.Repos{
			"web": fake.Branches{
				"master": fake.Branch{
					Meta: gh.Branch{
						SHA:     "bacdb99030b908fd853367f3c5cbbe20aa424672",
						Author:  "Hubot",
						Date:    time.Now(),
						Message: "fix the fix",
					},
					Files: fake.Files{
						{Location: &gh.Location{Path: ".deployment/ingress.libsonnet"}, Content: "github.com/rebuy-de/web/.deployment/ingress.libsonnet@master"},
						{Location: &gh.Location{Path: "blue/deployment/k8s/util.libsonnet"}, Content: "github.com/rebuy-de/web/blue/deployment/k8s/util.libsonnet@master"},
					},
				},
				"cloud-1337": fake.Branch{
					Meta: gh.Branch{
						SHA:     "2b8f82e0e1027b4c9f56da5a6175a099c9827859",
						Author:  "test",
						Date:    time.Now(),
						Message: "dis is test",
					},
					Files: fake.Files{
						{Location: &gh.Location{Path: ".deployment/ingress.libsonnet"}, Content: "github.com/rebuy-de/web/.deployment/ingress.libsonnet@cloud-1337"},
					},
				},
			},
			"jsonnet-libraries": fake.Branches{
				"master": fake.Branch{
					Meta: gh.Branch{
						SHA:     "15f088d8d78544cb4a900b28b8a2dcb28a130a3f",
						Author:  "root",
						Date:    time.Now(),
						Message: "initial",
					},
					Files: fake.Files{
						{Location: &gh.Location{Path: "util.libsonnet"}, Content: "github.com/rebuy-de/jsonnet-libraries/util.libsonnet@master"},
					},
				},
				"cloud-42": fake.Branch{
					Meta: gh.Branch{
						SHA:     "e95f454ee4b24573e6a38186d97edfa90a01f3c3",
						Author:  "user",
						Date:    time.Now(),
						Message: "add flux",
					},
					Files: fake.Files{
						{Location: &gh.Location{Path: "util.libsonnet"}, Content: "github.com/rebuy-de/jsonnet-libraries/util.libsonnet@cloud-42"},
					},
				},
			},
		},
	}

	importer := gh.NewJsonnetImporter(client)

	cases := []struct {
		name         string
		importedFrom string
		importedPath string
		foundAt      string
	}{
		{
			name:         "same_repo_no_ref_relative_path",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@cloud-1337",
			importedPath: "../../../.deployment/ingress.libsonnet",
			foundAt:      "github.com/rebuy-de/web/.deployment/ingress.libsonnet@cloud-1337",
		},
		{
			name:         "same_repo_no_ref_absolute_path",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@cloud-1337",
			importedPath: "github.com/rebuy-de/web/.deployment/ingress.libsonnet",
			foundAt:      "github.com/rebuy-de/web/.deployment/ingress.libsonnet@cloud-1337",
		},
		{
			name:         "same_repo_same_ref_relative_path",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@cloud-1337",
			importedPath: "../../../.deployment/ingress.libsonnet@cloud-1337",
			foundAt:      "github.com/rebuy-de/web/.deployment/ingress.libsonnet@cloud-1337",
		},
		{
			name:         "same_repo_other_ref_relative_path",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@cloud-1337",
			importedPath: "../../../.deployment/ingress.libsonnet@master",
			foundAt:      "github.com/rebuy-de/web/.deployment/ingress.libsonnet@master",
		},
		{
			name:         "other_repo_no_ref",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@cloud-1337",
			importedPath: "github.com/rebuy-de/jsonnet-libraries/util.libsonnet",
			foundAt:      "github.com/rebuy-de/jsonnet-libraries/util.libsonnet@master",
		},
		{
			name:         "other_repo_with_ref",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@cloud-1337",
			importedPath: "github.com/rebuy-de/jsonnet-libraries/util.libsonnet@cloud-42",
			foundAt:      "github.com/rebuy-de/jsonnet-libraries/util.libsonnet@cloud-42",
		},
		{
			name:         "local_import",
			importedFrom: "github.com/rebuy-de/web/blue/deployment/k8s/ingress.jsonnet@master",
			importedPath: "util.libsonnet",
			foundAt:      "github.com/rebuy-de/web/blue/deployment/k8s/util.libsonnet@master",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			contents, foundAt, err := importer.Import(tc.importedFrom, tc.importedPath)
			if err != nil {
				t.Error(err)
				return
			}

			if foundAt != tc.foundAt {
				t.Errorf("%s != %s", foundAt, tc.foundAt)
			}

			if contents.String() != tc.foundAt {
				t.Errorf("%s != %s", contents.String(), tc.foundAt)
			}
		})
	}

}
