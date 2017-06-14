package kubectl

import (
	"io"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/pkg/errors"
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

func (k *Kubectl) run(stdin io.Reader, args ...string) error {
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

	cmd.Stderr = log.WithFields(log.Fields{
		"executable": k.Path,
		"stream":     "stderr",
	}).WriterLevel(log.WarnLevel)
	cmd.Stdout = log.WithFields(log.Fields{
		"executable": k.Path,
		"stream":     "stdout",
	}).WriterLevel(log.DebugLevel)

	return cmd.Run()
}

func (k *Kubectl) Apply(stdin io.Reader) error {
	return k.run(stdin, "apply", "-f", "-")
}
