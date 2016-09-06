package kubernetes

import (
	"os/exec"
	"log"
	"strings"
	"github.com/rebuy-de/kubernetes-deployment/util"
	"io"
)

type Kubernetes struct {
	Kubeconfig  string
	KubectlPath string
}

func New(kubeconfig string) (*Kubernetes, error) {
	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return nil, err
	}

	return &Kubernetes{
		Kubeconfig:  kubeconfig,
		KubectlPath: kubectlPath,
	}, nil
}

func (k *Kubernetes) Exec(redirectStdOut bool, args ...string) (out []byte, err error) {
	var stderr , stdout io.ReadCloser

	cmd := exec.Command(k.KubectlPath, args...)
	log.Printf("$ kubectl %s", strings.Join(args, " "))

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return out, err
	}
	go util.PipeToLog("!", stderr)

	if redirectStdOut {
		stdout, err = cmd.StdoutPipe()
		if err != nil {
			return out, err
		}
		go util.PipeToLog(" ", stdout)
		err = cmd.Run()
	} else {
		out, err = cmd.Output()
		if err != nil {
			return out, err
		}

	}

	return out, err
}

func (k *Kubernetes) Apply(manifestFile string) error {
	_, err := k.Exec(true, "apply", "-f", manifestFile, "--record")
	return err
}

func (k *Kubernetes) Get(manifestFile string) ([]byte, error) {
	return k.Exec(false, "get", "-f", manifestFile, "-o", "json")
}

