package kubectl

import (
	"bytes"
	"encoding/json"
	"io"
	"os/exec"
	"strings"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/pkg/errors"
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

	stdout := new(bytes.Buffer)

	err = k.run(bytes.NewBuffer(raw), stdout, "apply", "-o", "json", "-f", "-")
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newObj, _, err := scheme.Codecs.UniversalDeserializer().Decode(stdout.Bytes(), nil, nil)
	return newObj, errors.Wrapf(err, "failed to decode result json")
}
