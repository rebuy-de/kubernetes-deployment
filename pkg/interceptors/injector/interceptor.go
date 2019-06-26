package injector

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
)

type Interceptor struct{}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return nil, err
	}

	annotations := accessor.GetAnnotations()
	if annotations == nil {
		return obj, nil
	}

	wantsInject := annotations["rebuy.com/kubernetes-deployment.inject-linkerd"]
	if wantsInject != "true" {
		return obj, nil
	}

	marshalled, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cmd := exec.Command(
		"linkerd", "inject",
		"--proxy-memory-request", "20Mi",
		"--proxy-cpu-request", "35m",
		"-")
	cmd.Stdin = bytes.NewBuffer(marshalled)
	cmd.Stderr = os.Stderr
	newUnmarshalled, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute linkerd")
	}

	newObj, err := kubeutil.Decode(newUnmarshalled)
	return newObj, errors.Wrapf(err, "failed to decode result json")
}
