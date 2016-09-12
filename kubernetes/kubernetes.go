package kubernetes

import (
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/rebuy-de/kubernetes-deployment/util"
)

type Kubernetes struct {
	Kubeconfig  string
	KubectlPath string
}

func New(kubeconfig string) (*Kubernetes, error) {
	kubectlPath, _ := exec.LookPath("kubectl")
	return &Kubernetes{
		Kubeconfig:  kubeconfig,
		KubectlPath: kubectlPath,
	}, nil
}

func (k *Kubernetes) Exec(args ...string) (out []byte, err error) {
	var stderr io.ReadCloser
	config := []string{}
	if k.Kubeconfig != "" {
		config = []string{"--kubeconfig=" + k.Kubeconfig}
	}
	config = append(config, args...)
	log.Printf("$ kubectl %s", strings.Join(config, " "))
	cmd := exec.Command(k.KubectlPath, config...)

	stderr, err = cmd.StderrPipe()
	if err != nil {
		return out, err
	}
	go util.PipeToLog("!", stderr)

	out, err = cmd.Output()
	if err != nil {
		return out, err
	}
	log.Printf("  %s", out)

	return out, err
}

func (k *Kubernetes) Apply(manifestFile string) ([]byte, error) {
	return k.Exec("apply", "-f", manifestFile, "--record")
}

func (k *Kubernetes) Get(manifestFile string) ([]byte, error) {
	return k.Exec("get", "-f", manifestFile, "-o", "json")
}
