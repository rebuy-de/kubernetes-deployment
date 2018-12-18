package injecter

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/rebuy-de/kubernetes-deployment/pkg/interceptors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	"github.com/rebuy-de/rebuy-go-sdk/testutil"
)

func TestTypePostManifestRenderer(t *testing.T) {
	var inter interceptors.PostManifestRenderer
	inter = New()
	_ = inter
}

func TestRendering(t *testing.T) {
	cases := []string{
		"deployment-inactive",
		"deployment",
	}

	for _, tc := range cases {
		t.Run(tc, func(t *testing.T) {
			srcFile := fmt.Sprintf("test-fixtures/%s.json", tc)
			dstFile := fmt.Sprintf("test-fixtures/%s-golden.json", tc)

			data, err := ioutil.ReadFile(srcFile)
			if err != nil {
				t.Error(err)
			}

			obj, err := kubeutil.Decode(data)
			if err != nil {
				t.Error(err)
			}

			newObj, err := New().PostManifestRender(obj)
			if err != nil {
				t.Error(err)
			}

			testutil.AssertGoldenJSON(t, dstFile, newObj)
		})
	}
}
