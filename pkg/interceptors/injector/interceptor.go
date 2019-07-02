package injector

import (
	"bytes"
	"encoding/json"
	"k8s.io/api/core/v1"
	"os"
	"os/exec"

	"github.com/pkg/errors"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
)

type Interceptor struct {
	Options Options
}

func New(options Options) *Interceptor {
	return &Interceptor{
		Options: options,
	}
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

	args := append([]string{"inject"}, append(i.Options.InjectArguments, "-")...)
	cmd := exec.Command("linkerd", args...)

	cmd.Stdin = bytes.NewBuffer(marshalled)
	cmd.Stderr = os.Stderr
	newUnmarshalled, err := cmd.Output()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to execute linkerd")
	}

	newObj, err := kubeutil.Decode(newUnmarshalled)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to decode result json")
	}

	template := kubeutil.PodTemplateSpecAccessor(newObj)
	if template != nil {
		for j, c := range template.Spec.Containers {
			if c.Name != "linkerd-proxy" {
				continue
			}

			template.Spec.Containers[j] = i.addLinkerdEnvVariables(c)
		}
	}

	return newObj, errors.Wrapf(err, "failed to decode result json")
}

func (i *Interceptor) addLinkerdEnvVariables(c v1.Container) v1.Container {
	c.Env = append(c.Env, v1.EnvVar{
		Name:  "LINKERD2_PROXY_OUTBOUND_CONNECT_TIMEOUT",
		Value: i.Options.ConnectTimeout,
	})

	return c
}
