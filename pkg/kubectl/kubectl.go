package kubectl

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
	log "github.com/sirupsen/logrus"
)

type Kubectl struct {
	Path       string
	Kubeconfig string
}

func New(path string, kubeconfig string) Interface {
	return &Kubectl{
		Path:       path,
		Kubeconfig: kubeconfig,
	}
}

func (k *Kubectl) run(stdin io.Reader, stdout io.Writer, args ...string) error {
	path, err := exec.LookPath(k.Path)
	if err != nil {
		return errors.Wrapf(err, "Unable to find kubectl executable '%s'", k.Path)
	}

	cmd := exec.Command(path)

	if strings.TrimSpace(k.Kubeconfig) != "" {
		cmd.Args = append(cmd.Args, "--kubeconfig", k.Kubeconfig)
	}

	cmd.Args = append(cmd.Args, args...)

	if stdin != nil {
		cmd.Stdin = stdin
	}

	if stdout != nil {
		cmd.Stdout = stdout
	}

	cmd.Stderr = log.WithFields(log.Fields{
		"executable": k.Path,
		"stream":     "stderr",
	}).WriterLevel(log.WarnLevel)

	return cmd.Run()
}

func (k *Kubectl) Apply(obj runtime.Object) (runtime.Object, error) {
	raw, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	stripped, err := emptyStatus(raw)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	stdout := new(bytes.Buffer)

	err = k.run(bytes.NewBuffer(stripped), stdout, "apply", "-o", "json", "-f", "-")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newObj, err := kubeutil.Decode(stdout.Bytes())
	return newObj, errors.Wrapf(err, "failed to decode result json")
}

func emptyStatus(raw []byte) ([]byte, error) {
	obj := make(map[string]interface{})
	if err := json.Unmarshal(raw, &obj); err != nil {
		return nil, err
	}
	if _, ok := obj["status"]; ok {
		obj["status"] = nil
	}

	return json.MarshalIndent(obj, "", "  ")
}
